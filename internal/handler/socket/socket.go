package socket

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"webrtc-server/driver"
	"webrtc-server/internal/repositories"
	"webrtc-server/internal/services"
	"webrtc-server/pkg/jwtauth"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Socket struct
type Socket struct {
	db   *driver.Database
	repo repositories.CallRepository
}

// InitSocket ...
func (s *Socket) InitSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	clientID := strconv.FormatInt(time.Now().UnixNano(), 10)

	client := &Client{
		ID:   clientID,
		hub:  hub,
		conn: conn,
		User: nil,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// NewSocketHandler handles websocket requests from the peer.
func NewSocketHandler(db *driver.Database) *Socket {
	return &Socket{
		db:   db,
		repo: services.NewCallService(db),
	}
}

// RegisterSocketID func
func (s *Socket) RegisterSocketID(token string, client *Client) {
	user, err := jwtauth.ParseTokenToUser(token, s.db)
	if err == nil {
		client.User = &user
		log.Println(client.User)
	} else {
		client.Emit("error", map[string]string{"error": "Invalid token, close connection"})
		client.Close()
	}
}

// InitSocketRoute register route for socket
func InitSocketRoute(socketHandler *Socket, routes *mux.Router) {
	socketHub := newHub()

	go socketHub.run(func(from *Client, target *Client, message *Message) {

		//register client
		if message.Action == "register" {
			socketHandler.RegisterSocketID(message.Data["token"], from)
			return
		}

		// Check client is registed
		if from.User == nil {
			from.Emit("error", map[string]string{"error": "Unregistered client"})
			from.Close()
			return
		}

		// handle action
		switch message.Action {
		case "call":
			if target == nil || from == nil {
				return
			}
			data := map[string]string{
				"from":     from.ID,
				"fullname": from.User.Fullname,
			}
			target.Emit("calling", data)
			break
		case "end_call":
			break
		}
	})

	routes.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socketHandler.InitSocket(socketHub, w, r)
	})
}

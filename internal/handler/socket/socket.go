package socket

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"webrtc-server/driver"
	"webrtc-server/internal/repositories"
	"webrtc-server/internal/services"
	"webrtc-server/pkg/helpers"

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

	clientID := strconv.FormatInt(time.Now().Unix(), 10)
	hash, err := helpers.HashAndSalt(clientID)

	client := &Client{
		ID:   hash,
		hub:  hub,
		conn: conn,
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
func (s *Socket) RegisterSocketID(token string) bool {
	return true
}

// InitSocketRoute register route for socket
func InitSocketRoute(socketHandler *Socket, routes *mux.Router) {
	socketHub := newHub()
	go socketHub.run(func(target *Client, message *Message) {
		switch message.Action {
		case "register":
			//success := socketHandler.RegisterSocketID(message.Data["token"])
			break
		case "call":
			data := map[string]string{
				"from": message.Data["token"],
			}
			target.Emit("calling", data)
			break
		case "end_call":
			break
		default:
		}
	})
	routes.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socketHandler.InitSocket(socketHub, w, r)
	})
}

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
	repo repositories.SocketRepository
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
		repo: services.NewSocketService(db),
	}
}

// RegisterSocketID func
func (s *Socket) RegisterSocketID(token string, client *Client) {
	user, err := jwtauth.ParseTokenToUser(token, s.db)
	if err == nil {
		client.User = &user
	} else {
		client.Emit("error", map[string]string{"error": "Invalid token, close connection"})
		client.Close()
	}
}

// MapEvents map events
func (s *Socket) MapEvents(from *Client, target *Client, message *Message) {

	//register client
	if message.Action == "register" {
		s.RegisterSocketID(message.Data["token"], from)
		return
	}

	// Check client is registed
	if from.User == nil {
		from.Emit("error", map[string]string{"error": "Unregistered client"})
		from.Close()
		return
	}

	if target == nil || from == nil {
		return
	}
	// handle action
	switch message.Action {
	case "call_start":
		// if target in another conversation
		if target.User.Calling {
			data := map[string]string{
				"message": "User is busy",
			}
			from.Emit("user_busy", data)
		}
		// emit to target
		data := map[string]string{
			"from":     from.ID,
			"fullname": from.User.Fullname,
		}
		target.Emit("call_start", data)
		return
	case "call_accepted":
		//set false for calling status
		target.User.Calling = true
		from.User.Calling = true
		s.repo.SetCallingStatus(true, []uint{target.User.ID, from.User.ID})

		//emit to target
		data := map[string]string{
			"message": from.User.Fullname + " accepted",
		}
		target.Emit("call_accepted", data)
		return
	case "call_busy":
		data := map[string]string{
			"message": from.User.Fullname + " is busy",
		}
		target.Emit("call_busy", data)
		return
	case "call_end":
		//set false for calling status
		target.User.Calling = false
		from.User.Calling = false
		s.repo.SetCallingStatus(false, []uint{target.User.ID, from.User.ID})
		//emit to target
		data := map[string]string{
			"message": from.User.Fullname + " ended the call",
		}

		target.Emit("call_end", data)
		return
	}
}

// InitSocketRoute register route for socket
func InitSocketRoute(s *Socket, routes *mux.Router) {
	socketHub := newHub()

	/**
	* Socket events
	* register, call, accept_call, end_call
	 */
	go socketHub.run(s.MapEvents)

	routes.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.InitSocket(socketHub, w, r)
	})
}

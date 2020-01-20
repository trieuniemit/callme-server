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
func (s *Socket) registerSocketID(token string, client *Client) {
	user, err := jwtauth.ParseTokenToUser(token, s.db)
	if err == nil {
		user.SocketID = client.ID
		client.User = &user

		socketIDs, err := s.repo.RegisterSocketID(&user)

		if err == nil {
			for _, ID := range socketIDs {
				if cl, ok := client.hub.clients[ID]; ok {
					cl.Emit("user_online", map[string]interface{}{
						"user": client.User,
					})
				}
			}
		} else {
			log.Println("Error when get contact socketIDs: ", err)
		}
		log.Println("Registed:", client.ID)
	} else {
		client.Emit("error", map[string]interface{}{"error": "Invalid token, close connection"})
		client.Close()
	}
}

func (s *Socket) setCallingStatus(status bool, from *Client, target *Client) {
	userIDs := []uint{}

	if target != nil {
		target.User.Calling = status
		userIDs = append(userIDs, target.User.ID)
	}
	if from != nil {
		from.User.Calling = status
		userIDs = append(userIDs, from.User.ID)
	}
	s.repo.SetCallingStatus(status, userIDs)
}

// MapEvents map events
func (s *Socket) MapEvents(from *Client, target *Client, message *Message) {

	if target != nil {
		log.Println(message.Action, ": ", from.ID, " - ", target.ID)
	}

	//register client
	if message.Action == "register" {
		s.registerSocketID(message.Data["token"].(string), from)
		return
	}
	// Check client is registed
	if from.User == nil {
		from.Emit("error", map[string]interface{}{"error": "Unregistered client"})
		from.Close()
		return
	}
	if target == nil && message.Action != "call_end" {
		from.Emit("call_not_available", map[string]interface{}{"error": "Tart get is not available"})
		return
	}
	if from == nil {
		return
	}
	// handle action
	switch message.Action {
	case "call_start":
		// if target in another conversation
		if target.User.Calling {
			data := map[string]interface{}{
				"message": "User is busy",
			}
			from.Emit("user_busy", data)
		}

		// emit to target
		s.setCallingStatus(true, from, target)
		data := map[string]interface{}{
			"user":        from.User,
			"description": message.Data["description"],
			"session_id":  message.Data["session_id"],
		}
		target.Emit("call_received", data)
		return
	case "call_accepted":
		//emit to target
		target.Emit("call_accepted", message.Data)
		return
	case "call_busy":
		data := map[string]interface{}{
			"message": from.User.Fullname + " is busy",
		}
		target.Emit("call_busy", data)
		return
	case "call_end":
		if target == nil {
			return
		}
		//set false for calling status
		s.setCallingStatus(false, from, target)
		//emit to target
		data := map[string]interface{}{
			"message": from.User.Fullname + " ended the call",
		}

		target.Emit("call_end", data)
		return
	case "call_candidate":
		if target == nil {
			return
		}
		//emit to target
		target.Emit(message.Action, message.Data)
		return
	case "call_answer":
		//emit to target
		target.Emit("call_accepted", message.Data)
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

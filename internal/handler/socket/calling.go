package socket

import (
	"log"
	"net/http"
	"webrtc-server/driver"

	"github.com/gorilla/websocket"

	"github.com/gorilla/mux"
)

// Socket is a middleman between the websocket connection and the server.
type Socket struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

// Message transfer to client
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// CallingWs struct
type CallingWs struct {
	db     *driver.Database
	socket *Socket
}

// WebSocketHandler of endpoint
func (callWs *CallingWs) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	cnn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		m := "Unable to upgrade to websockets"
		log.Println("conn err", err)
		http.Error(w, m, http.StatusBadRequest)
		return
	}
	// defer conn.Close()
	callWs.socket = &Socket{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}

	callWs.socket.clients[cnn] = true

	go callWs.handleReceiveMessages(cnn)
	go callWs.handleSendMessages(cnn)
}

func (callWs *CallingWs) handleSendMessages(nn *websocket.Conn) {
	for {
		// Grab the next message from the broadcast channel
		msg := <-callWs.socket.broadcast
		// Send it out to every client that is currently connected
		for client := range callWs.socket.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(callWs.socket.clients, client)
			}
		}
	}
}

func (callWs *CallingWs) handleReceiveMessages(cnn *websocket.Conn) {
	log.Println("WS: Connected========")
	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := cnn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(callWs.socket.clients, cnn)
			break
		}
		// Send the newly received message to the broadcast channel
		callWs.socket.broadcast <- msg
	}
}

// NewCallingWebSocket handler
func NewCallingWebSocket(db *driver.Database) *CallingWs {
	return &CallingWs{
		db: db,
	}
}

// RegisterWebSocket route
func RegisterWebSocket(wsHandler *CallingWs, routes *mux.Router) {
	routes.HandleFunc("/ws", wsHandler.WebSocketHandler)
}

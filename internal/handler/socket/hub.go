package socket

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) run(OnReceiveMessage func(*Client, *Message)) {
	log.Println("WS: init success =======")
	for {
		select {
		case client := <-h.register:
			log.Println("WS: new connection: ", client.ID)
			h.clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
				log.Println("WS: close connection: ", &client)
			}
		case message := <-h.broadcast:
			messageStrc := &Message{}

			// parse message to Message struct
			err := json.Unmarshal(message, messageStrc)
			if err == nil {
				// get target client
				if target, ok := h.clients[messageStrc.Target]; ok {
					OnReceiveMessage(target, messageStrc)
				}
			}

			// for _, client := range h.clients {
			// 	select {
			// 	case client.send <- message:
			// 		OnSendMessage(client.ID, msg)
			// 	default:
			// 		close(client.send)
			// 		delete(h.clients, client.ID)
			// 	}
			// }
		}
	}
}

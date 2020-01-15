package socket

import (
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

func (h *Hub) run(OnReceiveMessage func(*Client, *Client, *Message)) {
	log.Println("WS: init success =======")
	for {
		select {
		case client := <-h.register:
			log.Println("WS: new connection: ", client.ID)
			// var clientStr string
			// for clientID := range h.clients {
			// 	clientStr += clientID + "|"
			// }
			// client.send <- []byte(clientStr)
			h.clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
				log.Println("WS: close connection: ", &client)
			}
		case messageBytes := <-h.broadcast:
			message := MessageFromBytes(messageBytes)

			if message != nil {
				// get target client
				from := h.clients[message.Data["from"].(string)]
				targetID := message.Data["target"]
				if targetID != nil {
					target := h.clients[targetID.(string)]
					OnReceiveMessage(from, target, message)
				} else {
					OnReceiveMessage(from, nil, message)
				}

			}

		}
	}
}

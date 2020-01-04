package socket

import (
	"bytes"
	"encoding/json"
	"log"
)

// Message struct
type Message struct {
	Action string            `json:"action"`
	Data   map[string]string `json:"data"`
}

// ToBytes convert struct to []byte
func (m *Message) ToBytes() []byte {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(m)

	return reqBodyBytes.Bytes()
}

// MessageFromBytes func
func MessageFromBytes(msgBytes []byte) *Message {
	message := &Message{}

	// parse message to Message struct
	err := json.Unmarshal(msgBytes, message)
	if err == nil {
		// get target client
		return message
	}

	log.Println(err.Error())
	return nil
}

package socket

import (
	"bytes"
	"encoding/json"
)

// Message struct
type Message struct {
	Action  string            `json:"action"`
	Content map[string]string `json:"content"`
}

// ToBytes convert struct to []byte
func (m *Message) ToBytes() []byte {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(m)

	return reqBodyBytes.Bytes()
}

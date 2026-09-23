package websocket

type IncomingMessage struct {
	Type       string            `json:"type"`
	RoomID     string            `json:"room_id"`
	MessageID  string            `json:"message_id,omitempty"`
	Ciphertext string            `json:"ciphertext,omitempty"`
	Nonce      string            `json:"nonce,omitempty"`
	Keys       map[string]string `json:"keys,omitempty"`
}

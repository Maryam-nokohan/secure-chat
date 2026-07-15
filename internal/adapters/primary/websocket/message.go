package websocket

type IncomingMessage struct {
	Type       string            `json:"type"`
	RoomID     string            `json:"room_id"`
	Ciphertext string            `json:"ciphertext"`
	Nonce      string            `json:"nonce"`
	Keys       map[string]string `json:"keys"`
}
package message

const ChatSubject = "chat.messages"

type PubSubMessage struct {
	Type       string            `json:"type"`
	RoomID     string            `json:"room_id"`
	SenderID   string            `json:"sender_id"`
	Username   string            `json:"username"`
	Ciphertext string            `json:"ciphertext"`
	Nonce      string            `json:"nonce"`
	Keys       map[string]string `json:"keys"`
	Time       string            `json:"time"`
}
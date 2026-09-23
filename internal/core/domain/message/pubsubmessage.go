package message

const ChatSubject = "chat.messages"

type PubSubMessage struct {
	Type       string            `json:"type"`
	ID         string            `json:"id,omitempty"`
	MessageID  string            `json:"message_id,omitempty"`
	RoomID     string            `json:"room_id"`
	SenderID   string            `json:"sender_id"`
	Username   string            `json:"username"`
	Ciphertext string            `json:"ciphertext,omitempty"`
	Nonce      string            `json:"nonce,omitempty"`
	Keys       map[string]string `json:"keys,omitempty"`
	Time       string            `json:"time"`
}

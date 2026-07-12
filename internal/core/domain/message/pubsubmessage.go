package message

const ChatSubject = "chat.messages"

type PubSubMessage struct {
	Type     string `json:"type"`
	RoomID   string `json:"room_id"`
	SenderID string `json:"sender_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
	Time     string `json:"time"`
}
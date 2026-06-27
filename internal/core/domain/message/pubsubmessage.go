package message

type PubSubMessage struct {
    RoomID string `json:"room_id"`

    SenderID string `json:"sender_id"`

    Username string `json:"username"`

    Content string `json:"content"`

    CreatedAt string `json:"created_at"`
}
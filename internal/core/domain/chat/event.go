package chat

const PresenceSubject = "chat.presence"

type RoomEvent struct {
	Type     string `json:"type"`
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

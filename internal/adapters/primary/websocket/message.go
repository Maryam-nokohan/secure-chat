package websocket

type IncomingMessage struct {

    Type string `json:"type"`

    RoomID string `json:"room_id"`

    Content string `json:"content"`
}
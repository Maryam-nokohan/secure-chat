package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

type Client struct {
    ID       string
    Username string
    Conn     *websocket.Conn
    Send     chan []byte
    Room     string
    Hub      *Hub
}

type OutgoingMessage struct {
    Type     string `json:"type"`
    SenderID string `json:"sender_id"`
    Username string `json:"username"`
    RoomID   string `json:"room_id"`
    Content  string `json:"content"`
    Time     string `json:"time"`
}

func (c *Client) ReadPump(msgRepo ports.MessageRepository) {
    defer func() {
        c.Hub.Unregister <- c
        c.Conn.Close()
    }()

    for {
        _, raw, err := c.Conn.ReadMessage()
        if err != nil {
            break
        }

        var incoming IncomingMessage
        if err := json.Unmarshal(raw, &incoming); err != nil {
            continue
        }

        switch incoming.Type {
        case "join":
            c.Hub.JoinRoom(c, incoming.RoomID)

            ack, _ := json.Marshal(map[string]string{
                "type":    "joined",
                "room_id": incoming.RoomID,
            })
            select {
            case c.Send <- ack:
            default:
            }

        case "message":
            if c.Room == "" {
                continue
            }

            now := time.Now()

            msgID, _ := uuid.NewV4()
            go func() {
                _ = msgRepo.SaveMessage(context.Background(), &message.Message{
                    ID:             msgID,
                    RoomID:         uuid.FromStringOrNil(c.Room),
                    SenderID:       uuid.FromStringOrNil(c.ID),
                    SenderUsername: c.Username,
                    Content:        incoming.Content,
                    CreatedAt:      now,
                })
            }()

            out, err := json.Marshal(OutgoingMessage{
                Type:     "message",
                SenderID: c.ID,
                Username: c.Username,
                RoomID:   c.Room,
                Content:  incoming.Content,
                Time:     now.Format("15:04"),
            })
            if err != nil {
                continue
            }
            c.Hub.BroadcastToRoom(c.Room, out)
        }
    }
}

func (c *Client) WritePump() {
    defer c.Conn.Close()
    for msg := range c.Send {
        if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            return
        }
    }
    c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
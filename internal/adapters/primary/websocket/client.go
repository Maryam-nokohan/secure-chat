package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
	Room     string
	Hub      *Hub
	Broker   ports.MessageBroker
}

func (c *Client) ReadPump(msgSvc ports.MessageServiceI) {
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
				"type": "joined", "room_id": incoming.RoomID,
			})
			select {
			case c.Send <- ack:
			default:
			}

		case "message":
			if c.Room == "" {
				continue
			}

			roomID := uuid.FromStringOrNil(c.Room)
			senderID := uuid.FromStringOrNil(c.ID)

			saved, err := msgSvc.SaveGroupMessage(
				context.Background(), roomID, senderID, c.Username, incoming.Content,
			)
			if err != nil {
				continue
			}

			out := message.PubSubMessage{
				Type:     "message",
				SenderID: c.ID,
				Username: c.Username,
				RoomID:   c.Room,
				Content:  incoming.Content,
				Time:     saved.CreatedAt.Format(time.RFC3339),
			}
			payload, err := json.Marshal(out)
			if err != nil {
				continue
			}

			if pubErr := c.Broker.Publish(context.Background(), message.ChatSubject, payload); pubErr != nil {
				pkg.LogError(pubErr)
				c.Hub.BroadcastToRoom(c.Room, payload) 
			}
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
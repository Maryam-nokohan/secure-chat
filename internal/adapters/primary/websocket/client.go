// internal/adapters/primary/websocket/client.go
package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
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
	SenderID string `json:"sender_id"`
	Username string `json:"username"`
	RoomID   string `json:"room_id"`
	Content  string `json:"content"`
}

func (c *Client) ReadPump() {
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

		case "message":
			if c.Room == "" {
				continue
			}
			out, err := json.Marshal(OutgoingMessage{
				SenderID: c.ID,
				Username: c.Username,
				RoomID:   c.Room,
				Content:  incoming.Content,
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
	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
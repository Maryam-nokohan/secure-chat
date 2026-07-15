package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

type Handler struct {
	Hub     *Hub
	MsgSvc  ports.MessageServiceI
	Broker  ports.MessageBroker
	ChatSvc ports.ChatServiceI 
}

func NewHandler(hub *Hub, msgSvc ports.MessageServiceI, broker ports.MessageBroker, chatSvc ports.ChatServiceI) *Handler {
	return &Handler{Hub: hub, MsgSvc: msgSvc, Broker: broker, ChatSvc: chatSvc}
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true},
	}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	userID, existsID := c.Get("userID")
	username, existsName := c.Get("username")
	if !existsID || !existsName {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}

	client := &Client{
		Hub:      h.Hub,
		Broker:   h.Broker,
		ChatSvc:  h.ChatSvc,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		ID:       userID.(string),
		Username: username.(string),
	}
	h.Hub.Register <- client
	go client.WritePump()
	go client.ReadPump(h.MsgSvc)
}
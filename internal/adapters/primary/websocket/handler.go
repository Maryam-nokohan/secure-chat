package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	Hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		Hub: hub,
	}
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) ServeWs(w http.ResponseWriter, r *http.Request, username string) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	client := &Client{
		Id:       username,
		Username: username,
		Conn:     conn,
		Send:     make(chan []byte),
		Hub:      h.Hub,
	}

	h.Hub.Register <- client

	go client.writePump()
	go client.readPump()
}
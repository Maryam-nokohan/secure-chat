package websocket

import (
	"encoding/json"
	"sync"

	"github.com/maryam-nokohan/secure-chat/pkg"
)

type PresenceEvent struct {
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Online   bool   `json:"online"`
}

type Hub struct {
	mu         sync.Mutex
	Clients    map[string]*Client
	Rooms      map[string]map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[string]map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte, 256),
	}
}

func (h *Hub) Run() {
	pkg.LogInfo("WebSocket hub running...")
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()
			var memberOf []string
			if _, ok := h.Clients[client.ID]; ok {
				for roomID, room := range h.Rooms {
					if _, in := room[client.ID]; in {
						memberOf = append(memberOf, roomID)
					}
				}
				delete(h.Clients, client.ID)
				for _, room := range h.Rooms {
					delete(room, client.ID)
				}
				close(client.Send)
			}
			h.mu.Unlock()
			h.broadcastPresence(client, false, memberOf)

		case payload := <-h.Broadcast:
			h.mu.Lock()
			for id, client := range h.Clients {
				select {
				case client.Send <- payload:
				default:
					close(client.Send)
					delete(h.Clients, id)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) JoinRoom(client *Client, roomID string) {
    h.mu.Lock()
    if h.Rooms[roomID] == nil {
        h.Rooms[roomID] = make(map[string]*Client)
    }
    if client.Room != "" {
        delete(h.Rooms[client.Room], client.ID)
    }
    h.Rooms[roomID][client.ID] = client
    client.Room = roomID
    h.mu.Unlock()

    h.broadcastPresence(client, true, []string{roomID})
}

func (h *Hub) BroadcastToRoom(roomID string, payload []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, client := range h.Rooms[roomID] {
		select {
		case client.Send <- payload:
		default:
		}
	}
}

func (h *Hub) GetOnlineUserIDs() map[string]bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	online := make(map[string]bool, len(h.Clients))
	for id := range h.Clients {
		online[id] = true
	}
	return online
}

func (h *Hub) broadcastPresence(client *Client, online bool, roomIDs []string) {
    if len(roomIDs) == 0 {
        return
    }
    event, _ := json.Marshal(PresenceEvent{
        Type:     "presence",
        UserID:   client.ID,
        Username: client.Username,
        Online:   online,
    })
    h.mu.Lock()
    defer h.mu.Unlock()
    for _, roomID := range roomIDs {
        for cID, c := range h.Rooms[roomID] {
            if cID == client.ID {
                continue
            }
            select {
            case c.Send <- event:
            default:
            }
        }
    }
}

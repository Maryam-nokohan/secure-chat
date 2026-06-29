package websocket

import "sync"

type Hub struct {
	mu         sync.RWMutex
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
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				for _, room := range h.Rooms {
					delete(room, client.ID)
				}
				close(client.Send)
			}
			h.mu.Unlock()

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
	defer h.mu.Unlock()
	if h.Rooms[roomID] == nil {
		h.Rooms[roomID] = make(map[string]*Client)
	}
	if client.Room != "" {
		delete(h.Rooms[client.Room], client.ID)
	}
	h.Rooms[roomID][client.ID] = client
	client.Room = roomID
}

func (h *Hub) BroadcastToRoom(roomID string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if roomClients, ok := h.Rooms[roomID]; ok {
		for _, client := range roomClients {
			select {
			case client.Send <- payload:
			default:
			}
		}
	}
}
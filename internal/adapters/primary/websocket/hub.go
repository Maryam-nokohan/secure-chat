package websocket

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.Id] = client
			h.Broadcast <- []byte(client.Username + " joined the chat!")

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.Id]; ok {
				delete(h.Clients, client.Id)
				close(client.Send)
				h.Broadcast <- []byte(client.Username + " left the chat!")
			}

		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.Id)
				}
			}
		}
	}
}

func (h *Hub) joinRoom(client *Client, roomName string) {
	client.CurrentRoom = roomName
	h.Broadcast <- []byte(client.Username + " joined the room: " + roomName + " welcome!")
}

func (h *Hub) leaveRoom(client *Client) {
	if client.CurrentRoom != "" {
		h.Broadcast <- []byte(client.Username + " left the room: " + client.CurrentRoom + " bye!")
		client.CurrentRoom = ""
	}
}

func (h *Hub) sendMessageToRoom(client *Client, message string) {
	if client.CurrentRoom != "" {
		h.Broadcast <- []byte(client.Username + " in room " + client.CurrentRoom + ": " + message)
	} else {
		h.Broadcast <- []byte(client.Username + ": " + message)
	}
}
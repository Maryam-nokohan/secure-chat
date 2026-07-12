package chat

import (
	"context"
	"encoding/json"

	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/websocket"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type PubSubService struct {
	broker ports.MessageBroker
	hub    *websocket.Hub
}

func NewPubSubService(broker ports.MessageBroker, hub *websocket.Hub) *PubSubService {
	return &PubSubService{broker: broker, hub: hub}
}

func (p *PubSubService) Start(ctx context.Context) error {
	pkg.LogInfo("Subscribing to NATS subject: " + message.ChatSubject)

	return p.broker.Subscribe(ctx, message.ChatSubject, func(ctx context.Context ,payload []byte) {
		
		var event message.PubSubMessage

		if err := json.Unmarshal(payload, &event); err != nil {
			pkg.LogError(err)
			return 
		}
		p.hub.BroadcastToRoom(event.RoomID, payload)
	})
}
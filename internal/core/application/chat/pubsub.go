package chat

import (
	"context"
	"encoding/json"

	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/websocket"
	domainChat "github.com/maryam-nokohan/secure-chat/internal/core/domain/chat"
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
type routeEnvelope struct {
	RoomID string `json:"room_id"`
}

func (p *PubSubService) forward(ctx context.Context, payload []byte) {
	var env routeEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		pkg.LogError(err)
		return
	}
	p.hub.BroadcastToRoom(env.RoomID, payload)
}

func (p *PubSubService) Start(ctx context.Context) error {
	pkg.LogInfo("Subscribing to NATS subject: " + message.ChatSubject)
	if err := p.broker.Subscribe(ctx, message.ChatSubject, p.forward); err != nil {
		return err
	}
	pkg.LogInfo("Subscribing to NATS subject: " + domainChat.PresenceSubject)
	return p.broker.Subscribe(ctx, domainChat.PresenceSubject, p.forward)
}
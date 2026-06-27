package chat

import (
	"context"
	"encoding/json"

	"github.com/maryam-nokohan/secure-chat/internal/adapters/primary/websocket"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/redis"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
)

type PubSubService struct {
	redis *redis.Client
	hub   *websocket.Hub
}

func NewPubSubService(
	redis *redis.Client,
	hub *websocket.Hub,
) *PubSubService {
	return &PubSubService{
		redis: redis,
		hub:   hub,
	}
}

func (p *PubSubService) Start(
	ctx context.Context,
) {
	sub := p.redis.Subscribe(
		ctx,
		"chat_messages",
	)

	ch := sub.Channel()

	go func() {
		for msg := range ch {
			var event message.PubSubMessage
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				continue
			}
			raw, err := json.Marshal(event)
			if err != nil {
				continue
			}
			p.hub.BroadcastToRoom(event.RoomID, raw)
		}
	}()
}

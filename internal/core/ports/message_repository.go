package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, msg *message.Message) error
	GetMessagesByRoom(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*message.Message, error)
}
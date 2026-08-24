package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, msg *message.Message, keys []message.MessageKey) error
	GetMessageByID(ctx context.Context, msgID uuid.UUID) (*message.Message, error)
	GetRoomHistory(ctx context.Context, roomID uuid.UUID, limit int) ([]*message.Message, error)
	GetMessageKeysForUser(ctx context.Context, messageIDs []uuid.UUID, userID uuid.UUID) (map[uuid.UUID]string, error)
	DeleteMessage(ctx context.Context, msgID uuid.UUID) error
	EditMessage(ctx context.Context, msgID uuid.UUID, ciphertext, nonce string, keys []message.MessageKey) error
	CountMessages(ctx context.Context) (int64, error)
}
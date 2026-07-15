package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
)

type MessageServiceI interface {
	SaveGroupMessage(ctx context.Context, roomID, senderID uuid.UUID, senderUsername, ciphertext, nonce string, encryptedKeys map[string]string) (*message.Message, error)
	GetHistory(ctx context.Context, roomID, userID uuid.UUID) ([]*message.MessageWithKey, error)
	DeleteMessage(ctx context.Context, msgID, callerID uuid.UUID) error
	EditMessage(ctx context.Context, msgID, callerID uuid.UUID, ciphertext, nonce string, encryptedKeys map[string]string) error
}
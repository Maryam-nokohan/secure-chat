package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
)
type MessageServiceI interface {
    ProcessMessage(ctx context.Context, roomID, senderID uuid.UUID, encryptedContent string) (*message.Message, error)
    GetHistory(ctx context.Context, roomID uuid.UUID) ([]*message.Message, error)
}
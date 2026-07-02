package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
)

type MessageServiceI interface {
    SaveGroupMessage(ctx context.Context, roomID, senderID uuid.UUID, senderUsername, content string) (*message.Message, error)
    GetHistory(ctx context.Context, roomID uuid.UUID) ([]*message.Message, error)
    DeleteMessage(ctx context.Context, msgID, callerID uuid.UUID) error
    EditMessage(ctx context.Context, msgID, callerID uuid.UUID, content string) error
}
package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/chat"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
)

type ChatServiceI interface {
	CreateRoom(ctx context.Context, creatorID uuid.UUID, name string) (*chat.Room, error)
	JoinRoom(ctx context.Context, roomID, userID uuid.UUID) error
	LeaveRoom(ctx context.Context, roomId, userId uuid.UUID) error
}
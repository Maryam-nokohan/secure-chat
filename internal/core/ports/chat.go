package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/chat"
)

type ChatServiceI interface {
	CreateRoom(ctx context.Context, creatorID uuid.UUID) error
	JoinRoom(ctx context.Context, roomID, userID uuid.UUID) error
	LeaveRoom(ctx context.Context, roomId, userId uuid.UUID) error
	ListRooms(ctx context.Context )([]*chat.Room , error)
	SendMessage(roomId , userId uuid.UUID , content string) error
}

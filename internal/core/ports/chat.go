package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/chat"
)

type ChatServiceI interface {
	CreateRoom(ctx context.Context, creatorID uuid.UUID, name string) (*chat.Room, error)
	JoinRoom(ctx context.Context, roomID, userID uuid.UUID, username string) error
	JoinRoomByCode(ctx context.Context, code string, userID uuid.UUID, username string) (*chat.Room, error)
	LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error
	ListRooms(ctx context.Context) ([]*chat.Room, error)
	ListUserRooms(ctx context.Context, userID uuid.UUID) ([]*chat.Room, error)
	GetRoom(ctx context.Context, roomID uuid.UUID) (*chat.Room, error)
	MarkRoomRead(ctx context.Context, roomID, userID uuid.UUID) error
	GetUnreadRoomIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error)
}
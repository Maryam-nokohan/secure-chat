package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/chat"
)

type ChatRepositoryI interface {
	CreateRoom(ctx context.Context, room *chat.Room) error
	FindRoomByID(ctx context.Context, id uuid.UUID) (*chat.Room, error)
	FindRoomByName(ctx context.Context, name string) (*chat.Room, error)
	FindRoomByInviteCode(ctx context.Context, code string) (*chat.Room, error)
	ListRooms(ctx context.Context) ([]*chat.Room, error)
	ListUserRooms(ctx context.Context, userID uuid.UUID) ([]*chat.Room, error)
	DeleteRoom(ctx context.Context, id uuid.UUID) error
	AddUserToRoom(ctx context.Context, roomID, userID uuid.UUID) error
	RemoveUserFromRoom(ctx context.Context, roomID, userID uuid.UUID) error
	MarkRoomRead(ctx context.Context, roomID, userID uuid.UUID) error
	GetUnreadRoomIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error)
}

package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/chat"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
)

type ChatServiceI interface {
	CreateRoom(ctx context.Context, creatorID uuid.UUID) error
	JoinRoom(ctx context.Context, roomID, userID uuid.UUID) error
	LeaveRoom(ctx context.Context, roomId, userId uuid.UUID) error
	ListRooms(ctx context.Context )([]*chat.Room , error)
	SendMessage(roomId , userId uuid.UUID , content string) error
}

type MessageServiceI interface {
	SaveMessage(ctx context.Context, senderID, receiverID, roomID *uuid.UUID, text string) (*message.Message, error)
	GetMessagesByRoom(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*message.Message, error)
	GetMessagesBetweenUsers(ctx context.Context, userID1, userID2 uuid.UUID, limit, offset int) ([]*message.Message, error)
}
package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/redis"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/chat"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

type ChatService struct {
    chatRepo ports.ChatRepositoryI

    msgRepo ports.MessageRepository

    redis *redis.Client
}

func NewChatService(
    chatRepo ports.ChatRepositoryI,
    msgRepo ports.MessageRepository,
    redis *redis.Client,
) ports.ChatServiceI {

    return &ChatService{
        chatRepo: chatRepo,
        msgRepo: msgRepo,
        redis: redis,
    }
}
func (c *ChatService) CreateRoom(ctx context.Context, creatorID uuid.UUID, name string) (*chat.Room, error){
    if exist , _ := c.chatRepo.FindRoomByName(ctx ,name) ; exist != nil {
        return  nil , fmt.Errorf("This room %s already Exist.",name)
    }
    roomId , err := uuid.NewV4()
    if err != nil {
        return  nil , err
    }

    room := &chat.Room{
        ID: roomId,
        Name: name,
        CreatorID: creatorID,
        CreatedAt: time.Now() ,
    }
    c.chatRepo.CreateRoom(ctx , room)
    return  room , nil
}
func (c *ChatService) JoinRoom(ctx context.Context, roomID, userID uuid.UUID) error {
    _, err := c.chatRepo.FindRoomByID(ctx, roomID)
    if err != nil {
        return fmt.Errorf("room does not exist: %w", err)
    }

    if err := c.chatRepo.AddUserToRoom(ctx, roomID, userID); err != nil {
        return fmt.Errorf("failed to join room: %w", err)
    }

    return nil
}

func (c *ChatService) LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error {
    if err := c.chatRepo.RemoveUserFromRoom(ctx, roomID, userID); err != nil {
        return fmt.Errorf("failed to leave room: %w", err)
    }
    return nil
}

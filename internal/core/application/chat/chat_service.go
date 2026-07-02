package chat

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/adapters/secondary/redis"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/chat"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type ChatService struct {
	chatRepo ports.ChatRepositoryI
	msgRepo  ports.MessageRepository
	redis    *redis.Client
}

func NewChatService(
	chatRepo ports.ChatRepositoryI,
	msgRepo ports.MessageRepository,
	redis *redis.Client,
) ports.ChatServiceI {
	pkg.LogInfo("Init ChatService...")
	return &ChatService{chatRepo: chatRepo, msgRepo: msgRepo, redis: redis}
}

func generateInviteCode() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (c *ChatService) CreateRoom(ctx context.Context, creatorID uuid.UUID, name string) (*chat.Room, error) {

	existing, err := c.chatRepo.FindRoomByName(ctx, name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("room %q already exists", name)
	}

	inviteCode, err := generateInviteCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invite code: %w", err)
	}

	roomID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	room := &chat.Room{
		ID:         roomID,
		Name:       name,
		CreatorID:  creatorID,
		InviteCode: inviteCode,
		CreatedAt:  time.Now(),
	}
	if err := c.chatRepo.CreateRoom(ctx, room); err != nil {
		return nil, fmt.Errorf("create room: %w", err)
	}
	return room, nil
}

func (c *ChatService) JoinRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	room, err := c.chatRepo.FindRoomByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room does not exist: %w", err)
	}
	for _, u := range room.Users {
		if u.ID == userID {
			return nil
		}
	}
	if err := c.chatRepo.AddUserToRoom(ctx, roomID, userID); err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}
	return nil
}

func (c *ChatService) JoinRoomByCode(ctx context.Context, code string, userID uuid.UUID) (*chat.Room, error) {
	room, err := c.chatRepo.FindRoomByInviteCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("invalid invite code")
	}
	for _, u := range room.Users {
		if u.ID == userID {
			return room, nil
		}
	}
	if err := c.chatRepo.AddUserToRoom(ctx, room.ID, userID); err != nil {
		return nil, fmt.Errorf("failed to join room: %w", err)
	}
	return room, nil
}

func (c *ChatService) LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	if err := c.chatRepo.RemoveUserFromRoom(ctx, roomID, userID); err != nil {
		return fmt.Errorf("failed to leave room: %w", err)
	}
	return nil
}

func (c *ChatService) ListRooms(ctx context.Context) ([]*chat.Room, error) {
	return c.chatRepo.ListRooms(ctx)
}
func (c *ChatService) ListUserRooms(ctx context.Context, userID uuid.UUID) ([]*chat.Room, error) {
    return c.chatRepo.ListUserRooms(ctx, userID)
}
func (c *ChatService) GetRoom(ctx context.Context, roomID uuid.UUID) (*chat.Room, error) {
    return c.chatRepo.FindRoomByID(ctx, roomID)
}
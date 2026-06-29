package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/chat"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ports.ChatRepositoryI {
	pkg.LogInfo("Initializing ChatRepository...")
	return &ChatRepository{
		db: db,
	}
}

func (r *ChatRepository) CreateRoom(ctx context.Context, room *chat.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *ChatRepository) FindRoomByID(ctx context.Context, id uuid.UUID) (*chat.Room, error) {
	var room chat.Room
	if err := r.db.WithContext(ctx).Preload("Users").Where("id = ?", id).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *ChatRepository) FindRoomByName(ctx context.Context, name string) (*chat.Room, error) {
	var room chat.Room
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *ChatRepository) ListRooms(ctx context.Context) ([]*chat.Room, error) {
	var rooms []*chat.Room
	if err := r.db.WithContext(ctx).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *ChatRepository) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&chat.Room{}).Error
}

func (r *ChatRepository) AddUserToRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	room := chat.Room{ID: roomID}
	targetUser := user.User{ID: userID}

	return r.db.WithContext(ctx).Model(&room).Association("Users").Append(&targetUser)
}

func (r *ChatRepository) RemoveUserFromRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	room := chat.Room{ID: roomID}
	targetUser := user.User{ID: userID}

	return r.db.WithContext(ctx).Model(&room).Association("Users").Delete(&targetUser)
}

func (r *ChatRepository) FindRoomByInviteCode(ctx context.Context, code string) (*chat.Room, error) {
	var room chat.Room
	if err := r.db.WithContext(ctx).
		Preload("Users").
		Where("invite_code = ?", code).
		First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}
func (r *ChatRepository) ListUserRooms(ctx context.Context, userID uuid.UUID) ([]*chat.Room, error) {
    var rooms []*chat.Room
    err := r.db.WithContext(ctx).
        Joins("JOIN room_users ON room_users.room_id = rooms.id").
        Where("room_users.user_id = ?", userID).
        Find(&rooms).Error
    if err != nil {
        return nil, err
    }
    return rooms, nil
}
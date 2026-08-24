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
func (r *ChatRepository) MarkRoomRead(ctx context.Context, roomID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO room_reads (room_id, user_id, last_read_at)
		VALUES (?, ?, now())
		ON CONFLICT (room_id, user_id) DO UPDATE SET last_read_at = now()
	`, roomID, userID).Error
}

func (r *ChatRepository) GetUnreadRoomIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error) {
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).Raw(`
		SELECT r.id FROM rooms r
		JOIN room_users ru ON ru.room_id = r.id AND ru.user_id = ?
		LEFT JOIN room_reads rr ON rr.room_id = r.id AND rr.user_id = ?
		WHERE EXISTS (
			SELECT 1 FROM messages m
			WHERE m.room_id = r.id
			AND m.created_at > COALESCE(rr.last_read_at, 'epoch'::timestamptz)
		)
	`, userID, userID).Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	set := make(map[uuid.UUID]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return set, nil
}
func (r *ChatRepository) CountRooms(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&chat.Room{}).Count(&count).Error
	return count, err
}
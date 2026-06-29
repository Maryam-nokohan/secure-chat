package postgres

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"

	"gorm.io/gorm"
)

type MessageRepo struct {
	db *gorm.DB
}

func NewMessageRepository(
	db *gorm.DB,
) ports.MessageRepository {
	return &MessageRepo{
		db: db,
	}
}

func (r *MessageRepo) SaveMessage(ctx context.Context, msg *message.Message) error {
	if err := r.db.WithContext(ctx).Create(msg).Error; err != nil {
		return err
	}
	return nil
}

func (r *MessageRepo) GetRoomHistory(ctx context.Context, roomID uuid.UUID, limit int) ([]*message.Message, error) {
	var messages []*message.Message
	if err := r.db.WithContext(ctx).
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepo) DeleteMessage(ctx context.Context, msgID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&message.Message{}, "id = ?", msgID).Error
}

func (r *MessageRepo) EditMessage(ctx context.Context, msgID uuid.UUID, ciphertext, aesKey, iv string) error {
    return r.db.WithContext(ctx).
        Model(&message.Message{}).
        Where("id = ?", msgID).
        Updates(map[string]interface{}{
            "encrypted_content": ciphertext,
            "encrypted_aes_key": aesKey,
            "iv":                iv,
        }).Error
}
func (r * MessageRepo) GetMessageByID(ctx context.Context, msgID uuid.UUID) (*message.Message, error) {
	var msg message.Message
	if err := r.db.WithContext(ctx).First(&msg, "id = ?", msgID).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}
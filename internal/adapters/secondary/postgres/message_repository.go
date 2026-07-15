package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"gorm.io/gorm"
)

type MessageRepo struct{ db *gorm.DB }

func NewMessageRepository(db *gorm.DB) ports.MessageRepository {
	pkg.LogInfo("Initializing MessageRepository...")
	return &MessageRepo{db: db}
}

func (r *MessageRepo) SaveMessage(ctx context.Context, msg *message.Message, keys []message.MessageKey) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(msg).Error; err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := tx.Create(&keys).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *MessageRepo) GetRoomHistory(ctx context.Context, roomID uuid.UUID, limit int) ([]*message.Message, error) {
	var messages []*message.Message
	err := r.db.WithContext(ctx).
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *MessageRepo) GetMessageKeysForUser(ctx context.Context, messageIDs []uuid.UUID, userID uuid.UUID) (map[uuid.UUID]string, error) {
	if len(messageIDs) == 0 {
		return map[uuid.UUID]string{}, nil
	}
	var rows []message.MessageKey
	err := r.db.WithContext(ctx).
		Where("message_id IN ? AND recipient_id = ?", messageIDs, userID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[uuid.UUID]string, len(rows))
	for _, row := range rows {
		result[row.MessageID] = row.EncryptedKey
	}
	return result, nil
}

func (r *MessageRepo) DeleteMessage(ctx context.Context, msgID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&message.Message{}, "id = ?", msgID).Error
}

func (r *MessageRepo) EditMessage(ctx context.Context, msgID uuid.UUID, ciphertext, nonce string, keys []message.MessageKey) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&message.Message{}).Where("id = ?", msgID).
			Updates(map[string]interface{}{"content": ciphertext, "nonce": nonce}).Error; err != nil {
			return err
		}
		if err := tx.Where("message_id = ?", msgID).Delete(&message.MessageKey{}).Error; err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := tx.Create(&keys).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *MessageRepo) GetMessageByID(ctx context.Context, msgID uuid.UUID) (*message.Message, error) {
	var msg message.Message
	err := r.db.WithContext(ctx).First(&msg, "id = ?", msgID).Error
	return &msg, err
}
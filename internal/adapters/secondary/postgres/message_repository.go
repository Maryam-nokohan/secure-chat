package postgres

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(
	db *gorm.DB,
) ports.MessageRepository {

	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) SaveMessage(
	ctx context.Context,
	msg *message.Message,
) error {

	return r.db.
		WithContext(ctx).
		Create(msg).
		Error
}

func (r *MessageRepository) GetMessagesByRoom(
	ctx context.Context,
	roomID uuid.UUID,
	limit,
	offset int,
) ([]*message.Message, error) {

	var messages []*message.Message

	err := r.db.
		WithContext(ctx).
		Where("room_id = ?", roomID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&messages).
		Error

	return messages, err
}
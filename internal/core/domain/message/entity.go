package message

import (
	"time"

	"github.com/gofrs/uuid"
)


type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	RoomID         uuid.UUID `gorm:"type:uuid;index;not null"`
	SenderID       uuid.UUID `gorm:"type:uuid;index;not null"`
	SenderUsername string    `gorm:"type:varchar(50);not null"`
	Ciphertext     string    `gorm:"column:content;type:text;not null"`
	Nonce          string    `gorm:"column:nonce;type:varchar(64);not null"`
	CreatedAt      time.Time `gorm:"index"`
}


type MessageKey struct {
	MessageID    uuid.UUID `gorm:"type:uuid;primary_key;column:message_id"`
	RecipientID  uuid.UUID `gorm:"type:uuid;primary_key;column:recipient_id"`
	EncryptedKey string    `gorm:"type:text;not null;column:encrypted_key"`
}

func (MessageKey) TableName() string { return "message_keys" }

type MessageWithKey struct {
	*Message
	EncryptedKey string
}

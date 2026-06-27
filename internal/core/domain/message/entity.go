package message

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key"`
	RoomID           uuid.UUID `gorm:"type:uuid;index;not null"`
	SenderID         uuid.UUID `gorm:"type:uuid;index;not null"`
	RecipientID      uuid.UUID `gorm:"type:uuid;index;not null"`
	EncryptedContent string    `gorm:"type:text;not null"`
	EncryptedAESKey  string
	IV               string
	CreatedAt        time.Time `gorm:"index"`
}

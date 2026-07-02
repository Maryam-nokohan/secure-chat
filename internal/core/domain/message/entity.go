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
    Content        string    `gorm:"type:text;not null"`
    CreatedAt      time.Time `gorm:"index"`
}
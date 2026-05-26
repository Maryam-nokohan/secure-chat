package message

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
)

type Message struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SenderID   uuid.UUID  `gorm:"type:uuid;not null"`
	ReceiverID *uuid.UUID `gorm:"type:uuid"`
	RoomID     *uuid.UUID `gorm:"type:uuid"`
	Text       string     `gorm:"type:text;not null"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP"`

	Sender   user.User `gorm:"foreignKey:SenderID;references:ID"`
	Receiver user.User `gorm:"foreignKey:ReceiverID;references:ID"`
}

package chat

import (
	"github.com/gofrs/uuid"
	"time"
)

type Room struct {
	ID        uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string      `gorm:"size:100;not null"`
	CreatorID uuid.UUID   `gorm:"type:uuid;not null"`
	CreatedAt time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	Users     []uuid.UUID `gorm:"-:all"`
}

type RoomUser struct {
	RoomID   uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID   uuid.UUID `gorm:"type:uuid;primary_key"`
	JoinedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

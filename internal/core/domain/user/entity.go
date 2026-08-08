package user

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Username  string    `gorm:"uniqueIndex;size:50;not null"`
	PassHash  string    `gorm:"column:passhash;size:255;not null"`
	PublicKey string    `gorm:"type:text;not null"`
	Bio       string    `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
package user

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key"`
	Username  string         `gorm:"uniqueIndex;size:50;not null"`
	PassHash  string         `gorm:"column:passhash;size:255;not null"`
	PublicKey string         `gorm:"type:text;not null"`
	Bio       string         `gorm:"type:text"`
	Role      string         `gorm:"type:varchar(20);not null;default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
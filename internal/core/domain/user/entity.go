package user

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	Username      string    `gorm:"uniqueIndex;size:50;not null"`
	PassHash      string    `gorm:"column:passhash;size:255;not null"`
	PublicKey     string    `gorm:"type:text;not null"`
	Bio           string    `gorm:"type:text"`
	Email         string    `gorm:"type:varchar(255)"`
	AuthProvider  string    `gorm:"column:auth_provider;type:varchar(20);not null;default:'local'"`
	ProviderID    string    `gorm:"column:provider_id;type:varchar(255)"`
	EmailVerified bool      `gorm:"column:email_verified;not null;default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
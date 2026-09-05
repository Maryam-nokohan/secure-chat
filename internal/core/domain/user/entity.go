package user

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key"`
	Username  string         `gorm:"uniqueIndex;size:50;not null"`
	PassHash  string         `gorm:"column:passhash;size:255"`
	PublicKey string         `gorm:"type:text;not null"`

	WrappedPrivateKey string `gorm:"column:wrapped_private_key;type:text"`
	PrivateKeyIV      string `gorm:"column:private_key_iv;type:varchar(64)"`
	PrivateKeySalt    string `gorm:"column:private_key_salt;type:varchar(64)"`

	Email      string `gorm:"column:email;size:255;index"`
	Provider   string `gorm:"column:provider;size:20"` 
	ProviderID string `gorm:"column:provider_id;size:255;index"`

	Bio       string         `gorm:"type:text"`
	Role      string         `gorm:"type:varchar(20);not null;default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
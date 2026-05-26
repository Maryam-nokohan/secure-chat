package user

import "github.com/gofrs/uuid"

type User struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Username string    `gorm:"uniqueIndex;not null"`
	PassHash string    `gorm:"not null"`
	Bio     string
}

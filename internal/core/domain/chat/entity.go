package chat

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
)

type Room struct {
	ID        uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string       
	CreatorID uuid.UUID    
	CreatedAt time.Time   
	InviteCode string      `gorm:"uniqueIndex;size:32;not null"`
	Users     []user.User `gorm:"many2many:room_users;"`
}
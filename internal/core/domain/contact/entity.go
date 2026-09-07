package contact


import (
	"time"

	"github.com/gofrs/uuid"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusAccepted Status = "accepted"
	StatusDeclined Status = "declined"
	StatusBlocked  Status = "blocked"
)

type Contact struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	RequesterID uuid.UUID `gorm:"type:uuid;column:requester_id;index"`
	AddresseeID uuid.UUID `gorm:"type:uuid;column:addressee_id;index"`
	Status      Status    `gorm:"type:varchar(20);not null;default:'pending'"`
	RoomID      *uuid.UUID `gorm:"type:uuid;column:room_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Contact) TableName() string { return "contacts" }
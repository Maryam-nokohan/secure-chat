package attachment

import (
	"time"

	"github.com/gofrs/uuid"
)

type Attachment struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	RoomID      uuid.UUID `gorm:"type:uuid;index;not null"`
	UploaderID  uuid.UUID `gorm:"type:uuid;not null"`
	StorageKey  string    `gorm:"type:text;not null"`
	ContentType string    `gorm:"type:varchar(100);not null"`
	SizeBytes   int64     `gorm:"not null"`
	CreatedAt   time.Time
}

func (Attachment) TableName() string { return "attachments" }
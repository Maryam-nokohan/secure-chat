package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/attachment"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"gorm.io/gorm"
)

type AttachmentRepo struct{ db *gorm.DB }

func NewAttachmentRepository(db *gorm.DB) ports.AttachmentRepository {
	pkg.LogInfo("Initializing AttachmentRepository...")
	return &AttachmentRepo{db: db}
}

func (r *AttachmentRepo) Create(ctx context.Context, a *attachment.Attachment) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *AttachmentRepo) FindByID(ctx context.Context, id uuid.UUID) (*attachment.Attachment, error) {
	var a attachment.Attachment
	if err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
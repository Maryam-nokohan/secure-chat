package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/attachment"
)

type AttachmentRepository interface {
	Create(ctx context.Context, a *attachment.Attachment) error
	FindByID(ctx context.Context, id uuid.UUID) (*attachment.Attachment, error)
}
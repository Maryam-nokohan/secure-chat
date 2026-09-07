package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/contact"
)

type ContactRepository interface {
	Create(ctx context.Context, c *contact.Contact) error
	FindByID(ctx context.Context, id uuid.UUID) (*contact.Contact, error)
	FindByPair(ctx context.Context, a, b uuid.UUID) (*contact.Contact, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status contact.Status, roomID *uuid.UUID) error
	ListForUser(ctx context.Context, userID uuid.UUID, status contact.Status) ([]*contact.Contact, error)
}
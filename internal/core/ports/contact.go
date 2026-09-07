package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/contact"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
)

type ContactSummary struct {
	ID       string
	User     user.User
	Status   contact.Status
	Incoming bool
	RoomID   string
}

type ContactServiceI interface {
	SendRequest(ctx context.Context, requesterID uuid.UUID, targetPublicID string) (*contact.Contact, error)
	Accept(ctx context.Context, contactID, userID uuid.UUID) (*contact.Contact, error)
	Decline(ctx context.Context, contactID, userID uuid.UUID) error
	Block(ctx context.Context, contactID, userID uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]ContactSummary, error)
}
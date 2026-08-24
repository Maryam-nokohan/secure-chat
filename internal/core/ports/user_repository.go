package ports

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user user.User) error
	FindUserByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	FindUserByUsername(ctx context.Context, username string) (*user.User, error)
	EditUser(ctx context.Context, user user.User) error
	DeleteUser(ctx context.Context, user user.User) error

	ListUsers(ctx context.Context, limit, offset int) ([]user.User, error)
	CountUsers(ctx context.Context) (int64, error)
	RestoreUser(ctx context.Context, id uuid.UUID) error
}
package ports

import (
	"context"

	"github.com/maryam-nokohan/secure-chat/internal/core/domain/auth"
)

type UserServicesI interface {
	Register(ctx context.Context, username, password string) (*auth.AuthResult, error)
	Login(ctx context.Context, username, password string) (*auth.AuthResult, error)
}

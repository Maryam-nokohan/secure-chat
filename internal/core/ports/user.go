package ports

import (
	"context"

	"github.com/maryam-nokohan/secure-chat/internal/core/domain/auth"
)

type UserServicesI interface {
	Register(ctx context.Context, username, password, publicKey, wrappedPrivateKey, privateKeyIV, privateKeySalt string) (*auth.AuthResult, error)
	Login(ctx context.Context, username, password string) (*auth.AuthResult, error)

	FindOrCreateOAuthUser(ctx context.Context, info auth.UserInfo, provider string) (result *auth.AuthResult, isNew bool, needsKeys bool, err error)
	SetupEncryptionKeys(ctx context.Context, userID, publicKey, wrappedPrivateKey, privateKeyIV, privateKeySalt string) error
}

package ports

import "github.com/maryam-nokohan/secure-chat/internal/core/domain/auth"

type TokenService interface {
	Generate(userID, username, role string) (string, error)
	Validate(token string) (*auth.Claims, error)
}
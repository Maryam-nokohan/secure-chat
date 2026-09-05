package ports

import (
	"net/http"

	"github.com/maryam-nokohan/secure-chat/internal/core/domain/auth"
)

type OAuthService interface {
	BeginAuth(w http.ResponseWriter, r *http.Request, provider string)
	CompleteAuth(w http.ResponseWriter, r *http.Request, provider string) (*auth.UserInfo, error)
}
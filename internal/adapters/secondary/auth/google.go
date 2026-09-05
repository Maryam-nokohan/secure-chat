package auth

import (
	"errors"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	domainAuth "github.com/maryam-nokohan/secure-chat/internal/core/domain/auth"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

type GothOAuthService struct{}

func NewGothOAuthService(clientID, clientSecret, callbackURL string) ports.OAuthService {
	goth.UseProviders(
		google.New(clientID, clientSecret, callbackURL, "email", "profile"),
	)
	return &GothOAuthService{}
}

func (g *GothOAuthService) BeginAuth(w http.ResponseWriter, r *http.Request, provider string) {
	q := r.URL.Query()
	q.Set("provider", provider)
	r.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(w, r)
}

func (g *GothOAuthService) CompleteAuth(w http.ResponseWriter, r *http.Request, provider string) (*domainAuth.UserInfo, error) {
	q := r.URL.Query()
	q.Set("provider", provider)
	r.URL.RawQuery = q.Encode()

	gu, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return nil, err
	}
	if gu.UserID == "" {
		return nil, errors.New("oauth response missing subject id")
	}

	return &domainAuth.UserInfo{
		ProviderID:    gu.UserID,
		Email:         gu.Email,
		EmailVerified: gu.Email != "", 
		Name:          gu.Name,
	}, nil
}
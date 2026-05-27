package user

import (
	"context"
	"errors"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/auth"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type UserService struct {
	repo     ports.UserRepository
	tokenSvc ports.TokenService
}

func NewUSerService(repo ports.UserRepository , tokenSvc ports.TokenService) ports.UserServicesI {
	return &UserService{
		repo: repo,
		tokenSvc: tokenSvc,
	}
}

func (s *UserService) Register(ctx context.Context, name, password string) (*auth.AuthResult, error) {
	username := strings.TrimSpace(name)

	pkg.LogInfo("attempting user registration: " + username)

	existingUser, err := s.repo.FindUserByUsername(
		ctx,
		username,
	)

	if err == nil && existingUser != nil {

		pkg.LogInfo("registration failed, username already exists: " + username)
		return nil, errors.New("username already exists")
	}

	hash, err := pkg.HashPassword(password)
	if err != nil {

		pkg.LogError(err)
		return nil, err
	}

	userID, err := uuid.NewV4()
	if err != nil {

		pkg.LogError(err)
		return nil, err
	}

	newUser := user.User{
		ID:       userID,
		Username: username,
		PassHash: hash,
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {

		pkg.LogError(err)
		return nil, err
	}

	pkg.LogInfo("user registered successfully: " + username)

	token, err := s.tokenSvc.Generate(
		userID.String(),
		username,
	)

	if err != nil {
		pkg.LogError(err)
		return nil, err
	}

	return &auth.AuthResult{
		Token:    token,
		UserID:   userID.String(),
		Username: username,
	}, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*auth.AuthResult, error) {

	u, err := s.repo.FindUserByUsername(ctx, username)
	if err != nil {
		pkg.LogError(err)
		return nil, err
	}

	if err := pkg.CheckPassword(password, u.PassHash); err != nil {
		pkg.LogError(err)
		return nil, errors.New("invalid credentials")
	}

	token, err := s.tokenSvc.Generate(u.ID.String(), u.Username)
	if err != nil {
		pkg.LogError(err)
		return nil, err
	}

	return &auth.AuthResult{
		Token:    token,
		UserID:   u.ID.String(),
		Username: u.Username,
	}, nil

}

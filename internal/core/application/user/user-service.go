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

func NewUserService(
	repo ports.UserRepository,
	tokenSvc ports.TokenService,
) ports.UserServicesI {
	pkg.LogInfo("Init UserService...")

	return &UserService{
		repo:     repo,
		tokenSvc: tokenSvc,
	}
}

func (s *UserService) Register(
	ctx context.Context,
	name string,
	password string,
	publicKey string,
) (*auth.AuthResult, error) {
	username := strings.TrimSpace(name)

	if username == "" {
		return nil, errors.New("username is required")
	}

	if err := pkg.ValidatePassword(password); err != nil {
		return nil, err
	}

	if err := pkg.ValidateRSAPublicKey(publicKey); err != nil {
		return nil, err
	}

	existingUser, err := s.repo.FindUserByUsername(ctx, username)
	if err == nil && existingUser != nil {
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
		ID:        userID,
		Username:  username,
		PassHash:  hash,
		PublicKey: publicKey,
		Role:      "user",
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		pkg.LogError(err)
		return nil, err
	}

	token, err := s.tokenSvc.Generate(userID.String(), username, newUser.Role)
	if err != nil {
		pkg.LogError(err)
		return nil, err
	}

	pkg.LogInfo("User registered successfully: " + username)

	return &auth.AuthResult{
		Token:     token,
		UserID:    userID.String(),
		Username:  username,
		PublicKey: publicKey,
		Role:      newUser.Role,
	}, nil
}

func (s *UserService) Login(
	ctx context.Context,
	username string,
	password string,
) (*auth.AuthResult, error) {
	u, err := s.repo.FindUserByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := pkg.CheckPassword(password, u.PassHash); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.tokenSvc.Generate(u.ID.String(), u.Username, u.Role)
	if err != nil {
		pkg.LogError(err)
		return nil, err
	}

	return &auth.AuthResult{
		Token:     token,
		UserID:    u.ID.String(),
		Username:  u.Username,
		PublicKey: u.PublicKey,
		Role:      u.Role,
	}, nil
}

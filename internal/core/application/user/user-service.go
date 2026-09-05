package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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
	wrappedPrivateKey string,
	privateKeyIV string,
	privateKeySalt string,
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
	if wrappedPrivateKey == "" || privateKeyIV == "" || privateKeySalt == "" {
		return nil, errors.New("encryption key could not be generated on your device; please retry registration")
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
		ID:                userID,
		Username:          username,
		PassHash:          hash,
		PublicKey:         publicKey,
		WrappedPrivateKey: wrappedPrivateKey,
		PrivateKeyIV:      privateKeyIV,
		PrivateKeySalt:    privateKeySalt,
		Role:              "user",
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
		Token: token, UserID: userID.String(), Username: username,
		PublicKey: publicKey, Role: newUser.Role,
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

func (s *UserService) FindOrCreateOAuthUser(ctx context.Context, info auth.UserInfo, provider string) (*auth.AuthResult, bool, bool, error) {
	if info.ProviderID == "" {
		return nil, false, false, errors.New("oauth provider returned no subject id")
	}

	if existing, err := s.repo.FindUserByProvider(ctx, provider, info.ProviderID); err == nil && existing != nil {
		token, err := s.tokenSvc.Generate(existing.ID.String(), existing.Username, existing.Role)
		if err != nil {
			return nil, false, false, err
		}
		needsKeys := existing.PublicKey == "" || existing.WrappedPrivateKey == ""
		return &auth.AuthResult{
			Token: token, UserID: existing.ID.String(), Username: existing.Username,
			PublicKey: existing.PublicKey, Role: existing.Role,
		}, false, needsKeys, nil
	}
	if info.Email != "" {
	if byEmail, err := s.repo.FindUserByEmail(ctx, info.Email); err == nil && byEmail != nil {
		return nil, false, false, errors.New("an account with this email already exists; please log in with your original method")
	}
}

	username, err := s.uniqueUsernameFromEmail(ctx, info.Email, info.ProviderID)
	if err != nil {
		return nil, false, false, err
	}

	userID, err := uuid.NewV4()
	if err != nil {
		return nil, false, false, err
	}

	newUser := user.User{
		ID:         userID,
		Username:   username,
		Email:      info.Email,
		Provider:   provider,
		ProviderID: info.ProviderID,
		Role:       "user",
	}
	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		pkg.LogError(err)
		return nil, false, false, err
	}

	token, err := s.tokenSvc.Generate(userID.String(), username, newUser.Role)
	if err != nil {
		return nil, false, false, err
	}

	pkg.LogInfo("OAuth user provisioned: " + username + " via " + provider)

	return &auth.AuthResult{
		Token: token, UserID: userID.String(), Username: username, Role: newUser.Role,
	}, true, true, nil
}
func (s *UserService) SetupEncryptionKeys(ctx context.Context, userIDStr, publicKey, wrappedPrivateKey, privateKeyIV, privateKeySalt string) error {
	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		return errors.New("invalid user id")
	}
	u, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if u.PublicKey != "" {
		return errors.New("encryption keys already configured")
	}
	if err := pkg.ValidateRSAPublicKey(publicKey); err != nil {
		return err
	}
	if wrappedPrivateKey == "" || privateKeyIV == "" || privateKeySalt == "" {
		return errors.New("encryption key could not be generated on your device; please retry")
	}

	u.PublicKey = publicKey
	u.WrappedPrivateKey = wrappedPrivateKey
	u.PrivateKeyIV = privateKeyIV
	u.PrivateKeySalt = privateKeySalt

	return s.repo.EditUser(ctx, *u)
}

func (s *UserService) uniqueUsernameFromEmail(ctx context.Context, email, providerID string) (string, error) {
	base := email
	if at := strings.IndexByte(base, '@'); at > 0 {
		base = base[:at]
	}
	if base == "" {
		base = "user"
	}
	candidate := base
	for i := 0; i < 5; i++ {
		if _, err := s.repo.FindUserByUsername(ctx, candidate); err != nil {
			return candidate, nil
		}
		suffix := make([]byte, 3)
		if _, err := rand.Read(suffix); err != nil {
			return "", err
		}
		candidate = fmt.Sprintf("%s_%s", base, hex.EncodeToString(suffix))
	}
	return "", errors.New("could not allocate a unique username, please try again")
}

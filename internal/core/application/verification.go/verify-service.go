package verification

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

const (
	codeTTL        = 10 * time.Minute
	resendCooldown = 60 * time.Second
	maxAttempts    = 5
)

var (
	errInvalidCode  = errors.New("incorrect verification code")
	errCodeExpired  = errors.New("verification code expired or was never requested, request a new one")
	errTooManyTries = errors.New("too many wrong attempts, request a new code")
)

type Service struct {
	store    ports.VerificationStore
	mailer   ports.EmailSender
	userRepo ports.UserRepository
	key      []byte
}

func NewService(store ports.VerificationStore, mailer ports.EmailSender, userRepo ports.UserRepository, key []byte) ports.EmailVerificationServiceI {
	pkg.LogInfo("Init EmailVerificationService...")
	return &Service{store: store, mailer: mailer, userRepo: userRepo, key: key}
}

func (s *Service) hash(email, code string) string {
	m := hmac.New(sha256.New, s.key)
	m.Write([]byte("email-verify\x00" + email + "\x00" + code))
	return hex.EncodeToString(m.Sum(nil))
}

func (s *Service) SendCode(ctx context.Context, email string) error {
	email = pkg.NormalizeEmail(email)
	if err := pkg.ValidateEmail(email); err != nil {
		return err
	}
	if existing, err := s.userRepo.FindUserByEmail(ctx, email); err == nil && existing != nil {
		return errors.New("an account with this email already exists")
	}

	code, err := pkg.GenerateNumericCode(6)
	if err != nil {
		return err
	}
	if err := s.store.Put(ctx, email, s.hash(email, code), codeTTL, resendCooldown); err != nil {
		return err
	}
	if err := s.mailer.SendVerificationCode(ctx, email, code, codeTTL); err != nil {
		pkg.LogError(err)
		_ = s.store.Delete(ctx, email)
		return ports.ErrVerificationSendFailed
	}
	return nil
}

func (s *Service) VerifyCode(ctx context.Context, email, code string) error {
	email = pkg.NormalizeEmail(email)
	code = strings.TrimSpace(code)

	rec, err := s.store.Get(ctx, email)
	if errors.Is(err, ports.ErrVerificationNotFound) {
		return errCodeExpired
	}
	if err != nil {
		return err
	}
	if rec.Attempts >= maxAttempts {
		_ = s.store.Delete(ctx, email)
		return errTooManyTries
	}
	if !hmac.Equal([]byte(s.hash(email, code)), []byte(rec.CodeHash)) {
		n, err := s.store.RecordFailedAttempt(ctx, email)
		if err == nil && n >= maxAttempts {
			_ = s.store.Delete(ctx, email)
			return errTooManyTries
		}
		return errInvalidCode
	}
	return s.store.Delete(ctx, email)
}

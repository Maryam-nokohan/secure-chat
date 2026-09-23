package ports

import (
	"context"
	"errors"
	"time"
)

var (
	ErrVerificationNotFound   = errors.New("verification: no pending code")
	ErrVerificationCooldown   = errors.New("please wait a minute before requesting another code")
	ErrVerificationSendFailed = errors.New("could not send the verification email, please try again shortly")
)

type VerificationRecord struct {
	CodeHash string
	Attempts int
}

type VerificationStore interface {
	Put(ctx context.Context, email, codeHash string, ttl, resendCooldown time.Duration) error
	Get(ctx context.Context, email string) (*VerificationRecord, error)
	RecordFailedAttempt(ctx context.Context, email string) (int, error)
	Delete(ctx context.Context, email string) error
}

type EmailVerificationServiceI interface {
	SendCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email, code string) error 
}
package ports

import (
	"context"
	"time"
)

type EmailSender interface {
	SendVerificationCode(ctx context.Context, to, code string, ttl time.Duration) error
}
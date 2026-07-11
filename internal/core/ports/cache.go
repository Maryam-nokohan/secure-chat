package ports

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache: key not found")

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value any, tm time.Duration) error
	Delete(ctx context.Context, key string) error
}
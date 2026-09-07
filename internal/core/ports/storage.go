package ports

import (
	"context"
	"time"
)

type FileStorage interface {
	Upload(ctx context.Context, folder, filename string, content []byte, contentType string) (path string, err error)
	Download(ctx context.Context, path string) (content []byte, contentType string, err error)
	Delete(ctx context.Context, path string) error
	PresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error)

}
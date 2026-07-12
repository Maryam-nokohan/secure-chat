package redis

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type EncryptedCache struct {
	inner ports.Cache
	gcm   cipher.AEAD
}

func NewEncryptedCache(inner ports.Cache, key []byte) (ports.Cache, error) {
	if len(key) != 32 {
		return nil, errors.New("cache encryption key must be exactly 32 bytes (AES-256)")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	pkg.LogInfo("Cache encryption enabled (AES-256-GCM)")
	return &EncryptedCache{inner: inner, gcm: gcm}, nil
}

func (e *EncryptedCache) encrypt(plain []byte) (string, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := e.gcm.Seal(nonce, nonce, plain, nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func (e *EncryptedCache) decrypt(encoded string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	ns := e.gcm.NonceSize()
	if len(raw) < ns {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ct := raw[:ns], raw[ns:]
	return e.gcm.Open(nil, nonce, ct, nil)
}

func (e *EncryptedCache) Get(ctx context.Context, key string) ([]byte, error) {
	raw, err := e.inner.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var encoded string
	if err := json.Unmarshal(raw, &encoded); err != nil {
		return nil, err
	}
	return e.decrypt(encoded)
}

func (e *EncryptedCache) Set(ctx context.Context, key string, value any, tm time.Duration) error {
	plain, err := json.Marshal(value)
	if err != nil {
		return err
	}
	encoded, err := e.encrypt(plain)
	if err != nil {
		return err
	}
	return e.inner.Set(ctx, key, encoded, tm)
}

func (e *EncryptedCache) Delete(ctx context.Context, key string) error {
	return e.inner.Delete(ctx, key)
}

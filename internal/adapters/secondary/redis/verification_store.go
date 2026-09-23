package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/redis/go-redis/v9"
)

type VerificationStore struct{ rdb *redis.Client }

func NewVerificationStore(c *Client) ports.VerificationStore { return &VerificationStore{rdb: c.RDB} }

func codeKey(email string) string     { return "emailverify:code:" + email }
func cooldownKey(email string) string { return "emailverify:cd:" + email }

var incrIfExists = redis.NewScript(`
if redis.call('EXISTS', KEYS[1]) == 0 then return -1 end
return redis.call('HINCRBY', KEYS[1], 'attempts', 1)`)

func (s *VerificationStore) Put(ctx context.Context, email, codeHash string, ttl, cooldown time.Duration) error {
	ok, err := s.rdb.SetNX(ctx, cooldownKey(email), 1, cooldown).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ports.ErrVerificationCooldown
	}
	pipe := s.rdb.TxPipeline()
	pipe.Del(ctx, codeKey(email))
	pipe.HSet(ctx, codeKey(email), "hash", codeHash, "attempts", 0)
	pipe.Expire(ctx, codeKey(email), ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (s *VerificationStore) Get(ctx context.Context, email string) (*ports.VerificationRecord, error) {
	m, err := s.rdb.HGetAll(ctx, codeKey(email)).Result()
	if err != nil {
		return nil, err
	}
	if m["hash"] == "" {
		return nil, ports.ErrVerificationNotFound
	}
	attempts, _ := strconv.Atoi(m["attempts"])
	return &ports.VerificationRecord{CodeHash: m["hash"], Attempts: attempts}, nil
}

func (s *VerificationStore) RecordFailedAttempt(ctx context.Context, email string) (int, error) {
	n, err := incrIfExists.Run(ctx, s.rdb, []string{codeKey(email)}).Int()
	if err != nil {
		return 0, err
	}
	if n < 0 {
		return 0, ports.ErrVerificationNotFound
	}
	return n, nil
}

func (s *VerificationStore) Delete(ctx context.Context, email string) error {
	return s.rdb.Del(ctx, codeKey(email)).Err()
}
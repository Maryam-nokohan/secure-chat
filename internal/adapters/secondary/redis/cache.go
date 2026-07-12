package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *Client
}

func NewRedisCache(client *Client) ports.Cache {
	pkg.LogInfo("Initializing Redis Cache...")
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	pkg.LogInfo("Getting cache for key:" + key)
	val, err := r.client.RDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ports.ErrCacheMiss
		}
		return nil, err
	}
	return val, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, tm time.Duration) error {
	pkg.LogInfo("Setting cache for key:"+ key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.RDB.Set(ctx, key, data, tm).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.RDB.Del(ctx, key).Err()
}
package redis

import (
	"context"

	"github.com/maryam-nokohan/secure-chat/pkg"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *redis.Client
}

func NewRedis(addr string) *Client {
	pkg.LogInfo("Initializing Redis Client...")

	opts, err := redis.ParseURL(addr)
	if err != nil {
		pkg.LogFattal("invalid REDIS_ADDR: " + err.Error())
	}

	rdb := redis.NewClient(opts)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		pkg.LogFattal("failed to connect to redis: " + err.Error())
	}
	pkg.LogInfo("Redis connection established.")

	return &Client{RDB: rdb}
}
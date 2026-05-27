package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *goredis.Client
}

func NewRedis(addr string)*Client{
	rdb := goredis.NewClient(&goredis.Options{
		Addr: addr,
	})

	return &Client{
		RDB: rdb,
	}
}

func (c *Client) Publish(ctx context.Context , channel string , payload string) error {
	return c.RDB.Publish(ctx , channel , payload).Err()
}
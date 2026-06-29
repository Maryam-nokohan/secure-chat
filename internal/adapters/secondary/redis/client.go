package redis

import (
	"context"

	"github.com/maryam-nokohan/secure-chat/pkg"
	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *goredis.Client
}

func NewRedis(addr string)*Client{
	pkg.LogInfo("Initializing Redis Client...")
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
func (c *Client) Subscribe(
    ctx context.Context,
    channel string,
) *goredis.PubSub {

    return c.RDB.Subscribe(
        ctx,
        channel,
    )
}
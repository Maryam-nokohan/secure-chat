package nats

import (
	"context"
	"time"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"github.com/nats-io/nats.go"
)

type Client struct {
	conn *nats.Conn
}

func NewNATS(url string) (ports.MessageBroker , error){
	pkg.LogInfo("Inint NATS...")

	conn , err := nats.Connect(
		url,
		nats.Name("secure-chat"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func (_ *nats.Conn , err error)  {
			if err != nil {
				pkg.LogError(err)
			}
		}),
		nats.ReconnectHandler(func (_ *nats.Conn)  {
			pkg.LogInfo("Reconnecting to NATS...")
		}),
	)

	if err != nil {
		return nil , err
	}
	return &Client{
		conn: conn,
	} , nil
}

func (c *Client) Publish(ctx context.Context, subject string, payload []byte) error {
	return c.conn.Publish(subject, payload)
}
func (c *Client)Subscribe(ctx context.Context , subject string , handler func(ctx context.Context , payload []byte)) error{
	_, err := c.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(ctx ,msg.Data)
	})
	return err
}

func (c *Client) Close() {
	c.conn.Drain() 
}
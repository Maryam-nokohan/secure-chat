package nats

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"github.com/nats-io/nats.go"
)

const streamName = "CHAT_EVENTS"

type Client struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewNATS(url string) (ports.MessageBroker, error) {
	pkg.LogInfo("Initializing NATS (JetStream)...")

	conn, err := nats.Connect(
		url,
		nats.Name("secure-chat"),
		nats.Timeout(10*time.Second),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				pkg.LogError(err)
			}
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			pkg.LogInfo("Reconnected to NATS")
		}),
	)
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream(nats.MaxWait(10 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	if err := ensureStreamWithRetry(js); err != nil {
		return nil, fmt.Errorf("failed to ensure JetStream stream: %w", err)
	}

	pkg.LogInfo("NATS JetStream ready — stream " + streamName + " is persisting subjects chat.>")
	return &Client{conn: conn, js: js}, nil
}

func ensureStreamWithRetry(js nats.JetStreamContext) error {
	var lastErr error
	for i := 0; i < 10; i++ {
		if err := ensureStream(js); err == nil {
			return nil
		} else {
			lastErr = err
			pkg.LogInfo(fmt.Sprintf("JetStream not ready yet (attempt %d/10): %v", i+1, err))
			time.Sleep(2 * time.Second)
		}
	}
	return lastErr
}

func ensureStream(js nats.JetStreamContext) error {
	if _, err := js.StreamInfo(streamName); err == nil {
		return nil
	}
	_, err := js.AddStream(&nats.StreamConfig{
		Name:      streamName,
		Subjects:  []string{"chat.>"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		MaxAge:    24 * time.Hour,
	})
	return err
}

func (c *Client) Publish(ctx context.Context, subject string, payload []byte) error {
	_, err := c.js.Publish(subject, payload, nats.Context(ctx))
	return err
}

func (c *Client) Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, payload []byte)) error {
	durable := durableName(subject)

	_, err := c.js.Subscribe(subject, func(msg *nats.Msg) {
		defer func() {
			if r := recover(); r != nil {
				pkg.LogError(fmt.Errorf("panic handling NATS message on %s: %v", subject, r))
			}
		}()
		handler(ctx, msg.Data)
		if err := msg.Ack(); err != nil {
			pkg.LogError(err)
		}
	},
		nats.Durable(durable),
		nats.ManualAck(),
		nats.AckExplicit(),
		nats.DeliverNew(),
	)
	return err
}

func durableName(subject string) string {
	return "secure-chat-" + strings.ReplaceAll(subject, ".", "_")
}

func (c *Client) Close() {
	c.conn.Drain()
}

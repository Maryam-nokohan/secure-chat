package nats

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	streamName = "CHAT_EVENTS"

	// maxStreamBytes bounds disk usage independently of MaxAge, so a burst of
	// traffic can't fill the volume before the 24h age limit has a chance to prune it.
	maxStreamBytes = 1 << 30 // 1 GiB

	// ackWait must clear the slowest expected handler run (websocket fan-out to a
	// room's subscribers) so a healthy in-flight message is never redelivered
	// while still being processed.
	ackWait = 30 * time.Second

	// maxDeliver bounds redelivery attempts for a message whose handler keeps
	// failing, so a single poison message can't loop forever; it is Term'd after this.
	maxDeliver = 5

	// maxAckPending caps unacknowledged messages in flight per consumer, giving
	// JetStream backpressure instead of flooding the hub with more than it can process.
	maxAckPending = 512
)

type Client struct {
	conn *nats.Conn
	js   jetstream.JetStream

	mu        sync.Mutex
	consumers []jetstream.ConsumeContext
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

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	if err := ensureStreamWithRetry(js); err != nil {
		return nil, fmt.Errorf("failed to ensure JetStream stream: %w", err)
	}

	pkg.LogInfo("NATS JetStream ready — stream " + streamName + " is persisting subjects chat.>")
	return &Client{conn: conn, js: js}, nil
}

func ensureStreamWithRetry(js jetstream.JetStream) error {
	var lastErr error
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := ensureStream(ctx, js)
		cancel()
		if err == nil {
			return nil
		}
		lastErr = err
		pkg.LogInfo(fmt.Sprintf("JetStream not ready yet (attempt %d/10): %v", i+1, err))
		time.Sleep(2 * time.Second)
	}
	return lastErr
}

func ensureStream(ctx context.Context, js jetstream.JetStream) error {
	_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  []string{"chat.>"},
		Storage:   jetstream.FileStorage,
		Retention: jetstream.LimitsPolicy,
		MaxAge:    24 * time.Hour,
		MaxBytes:  maxStreamBytes,
	})
	return err
}

func (c *Client) Publish(ctx context.Context, subject string, payload []byte) error {
	_, err := c.js.Publish(ctx, subject, payload)
	return err
}

func (c *Client) Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, payload []byte)) error {
	durable := durableName(subject)

	stream, err := c.js.Stream(ctx, streamName)
	if err != nil {
		return fmt.Errorf("failed to look up stream %s: %w", streamName, err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       durable,
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckWait:       ackWait,
		MaxDeliver:    maxDeliver,
		MaxAckPending: maxAckPending,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer for %s: %w", subject, err)
	}

	consumeCtx, err := consumer.Consume(
		func(msg jetstream.Msg) {
			defer func() {
				if r := recover(); r != nil {
					pkg.LogError(fmt.Errorf("panic handling NATS message on %s: %v", subject, r))
					if termErr := msg.Term(); termErr != nil {
						pkg.LogError(termErr)
					}
					return
				}
				if ackErr := msg.Ack(); ackErr != nil {
					pkg.LogError(ackErr)
				}
			}()
			handler(ctx, msg.Data())
		},
		jetstream.ConsumeErrHandler(func(_ jetstream.ConsumeContext, err error) {
			pkg.LogError(fmt.Errorf("NATS consume error on %s: %w", subject, err))
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming %s: %w", subject, err)
	}

	c.mu.Lock()
	c.consumers = append(c.consumers, consumeCtx)
	c.mu.Unlock()

	go func() {
		<-ctx.Done()
		consumeCtx.Stop()
	}()

	return nil
}

func durableName(subject string) string {
	return "secure-chat-" + strings.ReplaceAll(subject, ".", "_")
}

func (c *Client) Close() {
	c.mu.Lock()
	consumers := c.consumers
	c.consumers = nil
	c.mu.Unlock()

	for _, cc := range consumers {
		cc.Drain()
		<-cc.Closed()
	}

	c.conn.Drain()
}

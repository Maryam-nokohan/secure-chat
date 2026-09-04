package nats

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/maryam-nokohan/secure-chat/pkg"
	"github.com/nats-io/nats.go"
)

// testNATSURL points at the JetStream instance provisioned by docker-compose
// (`docker compose up -d nats`). Override with TEST_NATS_URL to point elsewhere.
func testNATSURL() string {
	if url := os.Getenv("TEST_NATS_URL"); url != "" {
		return url
	}
	return "nats://127.0.0.1:4222"
}

// requireNATS skips the test fast when nothing is listening at url. NewNATS itself
// is built to retry for a long time against a slow-starting server (RetryOnFailedConnect
// plus ensureStreamWithRetry), so calling it directly against a dead address would make
// `go test ./...` hang for a long time on any machine without docker-compose running.
func requireNATS(t *testing.T, url string) {
	t.Helper()

	conn, err := nats.Connect(url, nats.Timeout(1*time.Second), nats.RetryOnFailedConnect(false))
	if err != nil {
		t.Skipf("NATS not reachable at %s, skipping integration test: %v", url, err)
	}
	conn.Close()
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	pkg.Init()

	url := testNATSURL()
	requireNATS(t, url)

	broker, err := NewNATS(url)
	if err != nil {
		t.Fatalf("NewNATS: %v", err)
	}

	client := broker.(*Client)
	t.Cleanup(client.Close)

	return client
}

// deleteTestConsumer removes the durable consumer created for subject so repeated
// test runs against a persistent JetStream volume don't accumulate orphaned consumers.
func deleteTestConsumer(t *testing.T, client *Client, subject string) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.js.DeleteConsumer(ctx, streamName, durableName(subject))
	})
}

func testSubject(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("chat.test.%s.%d", t.Name(), time.Now().UnixNano())
}

func TestClient_PublishSubscribe_RoundTrip(t *testing.T) {
	client := newTestClient(t)
	subject := testSubject(t)
	deleteTestConsumer(t, client, subject)

	received := make(chan []byte, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := client.Subscribe(ctx, subject, func(_ context.Context, payload []byte) {
		received <- payload
	}); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	// Let the pull consumer's Consume loop start fetching before publishing.
	time.Sleep(200 * time.Millisecond)

	want := []byte("hello-jetstream")
	if err := client.Publish(context.Background(), subject, want); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	select {
	case got := <-received:
		if string(got) != string(want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for message")
	}
}

func TestClient_Subscribe_MultipleMessagesDeliveredInOrder(t *testing.T) {
	client := newTestClient(t)
	subject := testSubject(t)
	deleteTestConsumer(t, client, subject)

	received := make(chan []byte, 3)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := client.Subscribe(ctx, subject, func(_ context.Context, payload []byte) {
		received <- payload
	}); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	for i := 0; i < 3; i++ {
		payload := []byte(fmt.Sprintf("msg-%d", i))
		if err := client.Publish(context.Background(), subject, payload); err != nil {
			t.Fatalf("Publish %d: %v", i, err)
		}
	}

	for i := 0; i < 3; i++ {
		select {
		case got := <-received:
			want := fmt.Sprintf("msg-%d", i)
			if string(got) != want {
				t.Fatalf("message %d: got %q, want %q", i, got, want)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out waiting for message %d", i)
		}
	}
}

func TestClient_Subscribe_StopsOnContextCancellation(t *testing.T) {
	client := newTestClient(t)
	subject := testSubject(t)
	deleteTestConsumer(t, client, subject)

	received := make(chan []byte, 2)
	ctx, cancel := context.WithCancel(context.Background())

	if err := client.Subscribe(ctx, subject, func(_ context.Context, payload []byte) {
		received <- payload
	}); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	if err := client.Publish(context.Background(), subject, []byte("before-cancel")); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	select {
	case <-received:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for first message")
	}

	cancel()
	time.Sleep(300 * time.Millisecond) // let the watcher goroutine call ConsumeContext.Stop()

	if err := client.Publish(context.Background(), subject, []byte("after-cancel")); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	select {
	case msg := <-received:
		t.Fatalf("received message after context cancellation: %s", msg)
	case <-time.After(500 * time.Millisecond):
	}
}

func TestClient_Close_DrainsWithoutHanging(t *testing.T) {
	client := newTestClient(t)
	subject := testSubject(t)
	deleteTestConsumer(t, client, subject)

	ctx := context.Background()
	if err := client.Subscribe(ctx, subject, func(context.Context, []byte) {}); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	done := make(chan struct{})
	go func() {
		client.Close()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Close did not return in time")
	}
}

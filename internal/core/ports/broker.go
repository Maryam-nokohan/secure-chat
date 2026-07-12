package ports

import "context"

type MessageBroker interface {
	Publish(ctx context.Context , subject string , payload []byte) error
	Subscribe(ctx context.Context , subject string , handler func(ctx context.Context , payload []byte)) error
	Close() 
}
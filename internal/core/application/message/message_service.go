package message

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

type MessageService struct {
	messageRepo ports.MessageRepository
}

func NewMessageService(msgRepo ports.MessageRepository) *MessageService {

	return &MessageService{
		messageRepo: msgRepo,
	}
}

func (s *MessageService) SavePrivateMessage(
	ctx context.Context,
	senderID uuid.UUID,
	receiverID uuid.UUID,
	content string,
) error {

	msg := &message.Message{
		ID:         uuid.Must(uuid.NewV4()),
		SenderID:   senderID,
		ReceiverID: &receiverID,
		Text:       content,
	}

	return s.messageRepo.SaveMessage(ctx, msg)
}

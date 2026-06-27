package message

import (
    "context"
    "time"

    "github.com/gofrs/uuid"
    "github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
    "github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

type MessageService struct {
    msgRepo ports.MessageRepository
}

func NewMessageService(repo ports.MessageRepository) ports.MessageServiceI {
    return &MessageService{
        msgRepo: repo,
    }
}

func (s *MessageService) ProcessMessage(ctx context.Context, roomID, senderID uuid.UUID, encryptedContent string) (*message.Message, error) {
    msgID, err := uuid.NewV4()
    if err != nil {
        return nil, err
    }

    msg := &message.Message{
        ID:               msgID,
        RoomID:           roomID,
        SenderID:         senderID,
        EncryptedContent: encryptedContent,
        CreatedAt:        time.Now(),
    }

    if err := s.msgRepo.SaveMessage(ctx, msg); err != nil {
        return nil, err
    }

    return msg, nil
}

func (s *MessageService) GetHistory(ctx context.Context, roomID uuid.UUID) ([]*message.Message, error) {
    return s.msgRepo.GetRoomHistory(ctx, roomID, 100)
}
func (s *MessageService) DeleteMessage(ctx context.Context, msgID uuid.UUID) error {
	return s.msgRepo.DeleteMessage(ctx, msgID)
}

func (s *MessageService) EditMessage(ctx context.Context, msgID uuid.UUID, rawMsg string) error {
	return s.msgRepo.EditMessage(ctx, msgID, newEncryptedContent)
}
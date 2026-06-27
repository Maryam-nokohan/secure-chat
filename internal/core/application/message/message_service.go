package message

import (
    "context"
    "time"
    "github.com/maryam-nokohan/secure-chat/internal/pkg"
    "github.com/gofrs/uuid"
    "github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
    "github.com/maryam-nokohan/secure-chat/internal/core/ports"
)

type MessageService struct {
    msgRepo ports.MessageRepository
    userRepo ports.UserRepository
    crypto *pkg.EncryptionService
}

func NewMessageService(repo ports.MessageRepository, userRepo ports.UserRepository, crypto *pkg.EncryptionService) ports.MessageServiceI {
    return &MessageService{
        msgRepo: repo,
        userRepo: userRepo,
        crypto: crypto,
    }
}

func (s *MessageService) ProcessMessage(ctx context.Context, roomID, senderID uuid.UUID, plaintext string) (*message.Message, error) {
    msgID, err := uuid.NewV4()
    if err != nil {
        return nil, err
    }

    recipient, err := s.userRepo.FindUserByID(ctx, senderID)
    if err != nil {
        return nil, err
    }
    
    encryptedContent, err := s.crypto.EncryptMessage(recipient.PublicKey, plaintext)

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

func (s *MessageService) EditMessage(ctx context.Context, msgID uuid.UUID, plaintext string) error {
    msg , err := s.msgRepo.GetMessageByID(ctx, msgID)
    if err != nil {
        return err
    }

    recipient, err := s.userRepo.FindUserByID(ctx, msg.SenderID)
    if err != nil {
        return err
    }

    newEncryptedContent, err := s.crypto.EncryptMessage(recipient.PublicKey, plaintext)
    if err != nil {
        return err
    }

    return s.msgRepo.EditMessage(ctx, msgID, newEncryptedContent)
}

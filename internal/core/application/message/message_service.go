package message

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type MessageService struct {
    msgRepo ports.MessageRepository
    userRepo ports.UserRepository
    crypto *pkg.HybridCryptoService
}

func NewMessageService(repo ports.MessageRepository, userRepo ports.UserRepository, crypto *pkg.HybridCryptoService) ports.MessageServiceI {
    return &MessageService{
        msgRepo: repo,
        userRepo: userRepo,
        crypto: crypto,
    }
}
func (s *MessageService) ProcessMessage(ctx context.Context, roomID, senderID, recipientID uuid.UUID, plaintext string) (*message.Message, error) {
	msgID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	recipient, err := s.userRepo.FindUserByID(ctx, recipientID)
	if err != nil {
		return nil, err
	}

	payload, err := s.crypto.Encrypt(recipient.PublicKey, []byte(plaintext))
	if err != nil {
		return nil, err
	}

	msg := &message.Message{
		ID:               msgID,
		RoomID:           roomID,
		SenderID:         senderID,
        RecipientID:      recipientID,
		EncryptedContent: payload.Ciphertext,
		EncryptedAESKey:  payload.EncryptedKey,
		IV:               payload.IV,
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
func (s *MessageService) DeleteMessage(ctx context.Context, msgID, callerID uuid.UUID) error {
    msg, err := s.msgRepo.GetMessageByID(ctx, msgID)
    if err != nil {
        return err
    }
    if msg.SenderID != callerID {
        return errors.New("forbidden: you did not send this message")
    }
    return s.msgRepo.DeleteMessage(ctx, msgID)
}

func (s *MessageService) EditMessage(ctx context.Context, msgID, callerID uuid.UUID, plaintext string) error {
    msg, err := s.msgRepo.GetMessageByID(ctx, msgID)
    if err != nil {
        return err
    }
    if msg.SenderID != callerID {
        return errors.New("forbidden: you did not send this message")
    }
    recipient, err := s.userRepo.FindUserByID(ctx, msg.RecipientID)
    if err != nil {
        return err
    }
    payload, err := s.crypto.Encrypt(recipient.PublicKey, []byte(plaintext))
    if err != nil {
        return err
    }
    return s.msgRepo.EditMessage(ctx, msgID, payload.Ciphertext)
}

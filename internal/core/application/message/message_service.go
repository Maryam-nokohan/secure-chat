package message

import (
    "context"
    "errors"
    "time"

    "github.com/gofrs/uuid"
    domainMsg "github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
    "github.com/maryam-nokohan/secure-chat/internal/core/ports"
    "github.com/maryam-nokohan/secure-chat/pkg"
)

type MessageService struct {
    msgRepo  ports.MessageRepository
    userRepo ports.UserRepository
}

func NewMessageService(repo ports.MessageRepository, userRepo ports.UserRepository) ports.MessageServiceI {
    pkg.LogInfo("Init MessageService...")
    return &MessageService{msgRepo: repo, userRepo: userRepo}
}

func (s *MessageService) SaveGroupMessage(ctx context.Context, roomID, senderID uuid.UUID, senderUsername, content string) (*domainMsg.Message, error) {
    msgID, err := uuid.NewV4()
    if err != nil {
        return nil, err
    }
    msg := &domainMsg.Message{
        ID:             msgID,
        RoomID:         roomID,
        SenderID:       senderID,
        SenderUsername: senderUsername,
        Content:        content,
        CreatedAt:      time.Now(),
    }
    if err := s.msgRepo.SaveMessage(ctx, msg); err != nil {
        return nil, err
    }
    return msg, nil
}

func (s *MessageService) GetHistory(ctx context.Context, roomID uuid.UUID) ([]*domainMsg.Message, error) {
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

func (s *MessageService) EditMessage(ctx context.Context, msgID, callerID uuid.UUID, content string) error {
    msg, err := s.msgRepo.GetMessageByID(ctx, msgID)
    if err != nil {
        return err
    }
    if msg.SenderID != callerID {
        return errors.New("forbidden: you did not send this message")
    }
    return s.msgRepo.EditMessage(ctx, msgID, content)
}
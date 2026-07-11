package message

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	domainMsg "github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

const historyCacheTTL = 5 * time.Minute

type MessageService struct {
    msgRepo  ports.MessageRepository
    userRepo ports.UserRepository
    cache    ports.Cache
}

func NewMessageService(repo ports.MessageRepository, userRepo ports.UserRepository, cache ports.Cache) ports.MessageServiceI {
    pkg.LogInfo("Init MessageService...")
    return &MessageService{msgRepo: repo, userRepo: userRepo, cache: cache}
}

func historyKey(roomID uuid.UUID) string { return "history:" + roomID.String() }

func (s *MessageService) SaveGroupMessage(ctx context.Context, roomID, senderID uuid.UUID, senderUsername, content string) (*domainMsg.Message, error) {
    msgID, err := uuid.NewV4()
    if err != nil {
        return nil, err
    }
    msg := &domainMsg.Message{
        ID: msgID, RoomID: roomID, SenderID: senderID,
        SenderUsername: senderUsername, Content: content, CreatedAt: time.Now(),
    }
    if err := s.msgRepo.SaveMessage(ctx, msg); err != nil {
        return nil, err
    }
    _ = s.cache.Delete(ctx, historyKey(roomID)) 
    return msg, nil
}

func (s *MessageService) GetHistory(ctx context.Context, roomID uuid.UUID) ([]*domainMsg.Message, error) {
    key := historyKey(roomID)

    if cached, err := s.cache.Get(ctx, key); err == nil {
        var messages []*domainMsg.Message
        if json.Unmarshal(cached, &messages) == nil {
            return messages, nil
        }
    }

    messages, err := s.msgRepo.GetRoomHistory(ctx, roomID, 100)
    if err != nil {
        return nil, err
    }
    _ = s.cache.Set(ctx, key, messages, historyCacheTTL)
    return messages, nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, msgID, callerID uuid.UUID) error {
    msg, err := s.msgRepo.GetMessageByID(ctx, msgID)
    if err != nil {
        return err
    }
    if msg.SenderID != callerID {
        return errors.New("forbidden: you did not send this message")
    }
    if err := s.msgRepo.DeleteMessage(ctx, msgID); err != nil {
        return err
    }
    _ = s.cache.Delete(ctx, historyKey(msg.RoomID))
    return nil
}

func (s *MessageService) EditMessage(ctx context.Context, msgID, callerID uuid.UUID, content string) error {
    msg, err := s.msgRepo.GetMessageByID(ctx, msgID)
    if err != nil {
        return err
    }
    if msg.SenderID != callerID {
        return errors.New("forbidden: you did not send this message")
    }
    if err := s.msgRepo.EditMessage(ctx, msgID, content); err != nil {
        return err
    }
    _ = s.cache.Delete(ctx, historyKey(msg.RoomID))
    return nil
}
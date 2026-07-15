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

func keysToRows(msgID uuid.UUID, encryptedKeys map[string]string) []domainMsg.MessageKey {
	rows := make([]domainMsg.MessageKey, 0, len(encryptedKeys))
	for uidStr, encKey := range encryptedKeys {
		uid, err := uuid.FromString(uidStr)
		if err != nil {
			continue
		}
		rows = append(rows, domainMsg.MessageKey{MessageID: msgID, RecipientID: uid, EncryptedKey: encKey})
	}
	return rows
}

func (s *MessageService) SaveGroupMessage(ctx context.Context, roomID, senderID uuid.UUID, senderUsername, ciphertext, nonce string, encryptedKeys map[string]string) (*domainMsg.Message, error) {
	if ciphertext == "" || nonce == "" || len(encryptedKeys) == 0 {
		return nil, errors.New("message must be encrypted with at least one recipient key")
	}

	msgID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	msg := &domainMsg.Message{
		ID: msgID, RoomID: roomID, SenderID: senderID,
		SenderUsername: senderUsername, Ciphertext: ciphertext, Nonce: nonce, CreatedAt: time.Now(),
	}

	if err := s.msgRepo.SaveMessage(ctx, msg, keysToRows(msgID, encryptedKeys)); err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, historyKey(roomID))
	return msg, nil
}

func (s *MessageService) GetHistory(ctx context.Context, roomID, userID uuid.UUID) ([]*domainMsg.MessageWithKey, error) {
	key := historyKey(roomID)
	var messages []*domainMsg.Message

	if cached, err := s.cache.Get(ctx, key); err == nil {
		_ = json.Unmarshal(cached, &messages)
	}
	if messages == nil {
		var err error
		messages, err = s.msgRepo.GetRoomHistory(ctx, roomID, 100)
		if err != nil {
			return nil, err
		}
		_ = s.cache.Set(ctx, key, messages, historyCacheTTL)
	}

	ids := make([]uuid.UUID, len(messages))
	for i, m := range messages {
		ids[i] = m.ID
	}
	keyMap, err := s.msgRepo.GetMessageKeysForUser(ctx, ids, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*domainMsg.MessageWithKey, 0, len(messages))
	for _, m := range messages {
		if encKey, ok := keyMap[m.ID]; ok {
			result = append(result, &domainMsg.MessageWithKey{Message: m, EncryptedKey: encKey})
		}
	}
	return result, nil
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

func (s *MessageService) EditMessage(ctx context.Context, msgID, callerID uuid.UUID, ciphertext, nonce string, encryptedKeys map[string]string) error {
	if ciphertext == "" || nonce == "" || len(encryptedKeys) == 0 {
		return errors.New("edited message must be encrypted with at least one recipient key")
	}
	msg, err := s.msgRepo.GetMessageByID(ctx, msgID)
	if err != nil {
		return err
	}
	if msg.SenderID != callerID {
		return errors.New("forbidden: you did not send this message")
	}
	if err := s.msgRepo.EditMessage(ctx, msgID, ciphertext, nonce, keysToRows(msgID, encryptedKeys)); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, historyKey(msg.RoomID))
	return nil
}
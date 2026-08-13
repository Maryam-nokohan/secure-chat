package message

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	domainMsg "github.com/maryam-nokohan/secure-chat/internal/core/domain/message"
)

type mockMessageRepository struct {
	saveMessageFn           func(context.Context, *domainMsg.Message, []domainMsg.MessageKey) error
	getMessageByIDFn        func(context.Context, uuid.UUID) (*domainMsg.Message, error)
	getRoomHistoryFn        func(context.Context, uuid.UUID, int) ([]*domainMsg.Message, error)
	getMessageKeysForUserFn func(context.Context, []uuid.UUID, uuid.UUID) (map[uuid.UUID]string, error)
	deleteMessageFn         func(context.Context, uuid.UUID) error
	editMessageFn           func(context.Context, uuid.UUID, string, string, []domainMsg.MessageKey) error

	saveCalls       int
	getMessageCalls int
	historyCalls    int
	keysCalls       int
	deleteCalls     int
	editCalls       int

	savedMessage *domainMsg.Message
	savedKeys    []domainMsg.MessageKey

	lastHistoryRoomID uuid.UUID
	lastHistoryLimit  int

	lastKeysMessageIDs []uuid.UUID
	lastKeysUserID     uuid.UUID

	lastDeletedMessageID uuid.UUID

	lastEditedMessageID uuid.UUID
	lastCiphertext      string
	lastNonce           string
	lastEditedKeys      []domainMsg.MessageKey
}

func (m *mockMessageRepository) SaveMessage(
	ctx context.Context,
	msg *domainMsg.Message,
	keys []domainMsg.MessageKey,
) error {
	m.saveCalls++

	m.savedMessage = msg
	m.savedKeys = append([]domainMsg.MessageKey(nil), keys...)

	if m.saveMessageFn != nil {
		return m.saveMessageFn(ctx, msg, keys)
	}

	return nil
}

func (m *mockMessageRepository) GetMessageByID(
	ctx context.Context,
	msgID uuid.UUID,
) (*domainMsg.Message, error) {
	m.getMessageCalls++

	if m.getMessageByIDFn != nil {
		return m.getMessageByIDFn(ctx, msgID)
	}

	return nil, errors.New("GetMessageByID not configured")
}

func (m *mockMessageRepository) GetRoomHistory(
	ctx context.Context,
	roomID uuid.UUID,
	limit int,
) ([]*domainMsg.Message, error) {
	m.historyCalls++

	m.lastHistoryRoomID = roomID
	m.lastHistoryLimit = limit

	if m.getRoomHistoryFn != nil {
		return m.getRoomHistoryFn(ctx, roomID, limit)
	}

	return nil, nil
}

func (m *mockMessageRepository) GetMessageKeysForUser(
	ctx context.Context,
	messageIDs []uuid.UUID,
	userID uuid.UUID,
) (map[uuid.UUID]string, error) {
	m.keysCalls++

	m.lastKeysMessageIDs = append([]uuid.UUID(nil), messageIDs...)
	m.lastKeysUserID = userID

	if m.getMessageKeysForUserFn != nil {
		return m.getMessageKeysForUserFn(ctx, messageIDs, userID)
	}

	return map[uuid.UUID]string{}, nil
}

func (m *mockMessageRepository) DeleteMessage(
	ctx context.Context,
	msgID uuid.UUID,
) error {
	m.deleteCalls++

	m.lastDeletedMessageID = msgID

	if m.deleteMessageFn != nil {
		return m.deleteMessageFn(ctx, msgID)
	}

	return nil
}

func (m *mockMessageRepository) EditMessage(
	ctx context.Context,
	msgID uuid.UUID,
	ciphertext string,
	nonce string,
	keys []domainMsg.MessageKey,
) error {
	m.editCalls++

	m.lastEditedMessageID = msgID
	m.lastCiphertext = ciphertext
	m.lastNonce = nonce
	m.lastEditedKeys = append([]domainMsg.MessageKey(nil), keys...)

	if m.editMessageFn != nil {
		return m.editMessageFn(ctx, msgID, ciphertext, nonce, keys)
	}

	return nil
}

type mockCache struct {
	values map[string][]byte

	getCalls    int
	setCalls    int
	deleteCalls int

	lastGetKey    string
	lastSetKey    string
	lastSetValue  any
	lastSetTTL    time.Duration
	lastDeleteKey string

	getErr    error
	setErr    error
	deleteErr error
}

func newMockCache() *mockCache {
	return &mockCache{
		values: make(map[string][]byte),
	}
}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.getCalls++
	m.lastGetKey = key

	if m.getErr != nil {
		return nil, m.getErr
	}

	value, ok := m.values[key]
	if !ok {
		return nil, errors.New("cache miss")
	}

	return append([]byte(nil), value...), nil
}

func (m *mockCache) Set(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) error {
	m.setCalls++
	m.lastSetKey = key
	m.lastSetValue = value
	m.lastSetTTL = ttl

	if m.setErr != nil {
		return m.setErr
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	m.values[key] = data

	return nil
}

func (m *mockCache) Delete(ctx context.Context, key string) error {
	m.deleteCalls++
	m.lastDeleteKey = key

	if m.deleteErr != nil {
		return m.deleteErr
	}

	delete(m.values, key)

	return nil
}

func newTestMessage(
	id uuid.UUID,
	roomID uuid.UUID,
	senderID uuid.UUID,
	ciphertext string,
	nonce string,
) *domainMsg.Message {
	return &domainMsg.Message{
		ID:             id,
		RoomID:         roomID,
		SenderID:       senderID,
		SenderUsername: "alice",
		Ciphertext:     ciphertext,
		Nonce:          nonce,
		CreatedAt:      time.Now(),
	}
}

func newTestService(
	repo *mockMessageRepository,
	cache *mockCache,
) *MessageService {
	return &MessageService{
		msgRepo: repo,
		cache:   cache,
	}
}

func TestSaveGroupMessage_Success(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	service := newTestService(repo, cache)

	roomID := uuid.Must(uuid.NewV4())
	senderID := uuid.Must(uuid.NewV4())
	recipientID := uuid.Must(uuid.NewV4())

	keys := map[string]string{
		recipientID.String(): "encrypted-key",
	}

	msg, err := service.SaveGroupMessage(
		context.Background(),
		roomID,
		senderID,
		"alice",
		"ciphertext",
		"nonce",
		keys,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if msg == nil {
		t.Fatal("expected message, got nil")
	}

	if msg.ID == uuid.Nil {
		t.Fatal("expected generated message ID")
	}

	if msg.RoomID != roomID {
		t.Fatalf("expected room ID %s, got %s", roomID, msg.RoomID)
	}

	if msg.SenderID != senderID {
		t.Fatalf("expected sender ID %s, got %s", senderID, msg.SenderID)
	}

	if msg.SenderUsername != "alice" {
		t.Fatalf("expected username alice, got %s", msg.SenderUsername)
	}

	if msg.Ciphertext != "ciphertext" {
		t.Fatalf("expected ciphertext to be preserved")
	}

	if msg.Nonce != "nonce" {
		t.Fatalf("expected nonce to be preserved")
	}

	if repo.saveCalls != 1 {
		t.Fatalf("expected 1 SaveMessage call, got %d", repo.saveCalls)
	}

	if len(repo.savedKeys) != 1 {
		t.Fatalf("expected 1 message key, got %d", len(repo.savedKeys))
	}

	if repo.savedKeys[0].MessageID != msg.ID {
		t.Fatalf("expected key message ID %s, got %s",
			msg.ID,
			repo.savedKeys[0].MessageID,
		)
	}

	if repo.savedKeys[0].RecipientID != recipientID {
		t.Fatalf("expected recipient ID %s, got %s",
			recipientID,
			repo.savedKeys[0].RecipientID,
		)
	}

	if repo.savedKeys[0].EncryptedKey != "encrypted-key" {
		t.Fatalf("expected encrypted key to be preserved")
	}

	if cache.deleteCalls != 1 {
		t.Fatalf("expected cache invalidation, got %d calls", cache.deleteCalls)
	}

	expectedCacheKey := "history:" + roomID.String()

	if cache.lastDeleteKey != expectedCacheKey {
		t.Fatalf(
			"expected cache key %q, got %q",
			expectedCacheKey,
			cache.lastDeleteKey,
		)
	}
}

func TestSaveGroupMessage_RejectsEmptyCiphertext(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	service := newTestService(repo, cache)

	_, err := service.SaveGroupMessage(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
		"alice",
		"",
		"nonce",
		map[string]string{
			uuid.Must(uuid.NewV4()).String(): "key",
		},
	)

	if err == nil {
		t.Fatal("expected error for empty ciphertext")
	}

	if repo.saveCalls != 0 {
		t.Fatalf("repository should not be called, got %d calls", repo.saveCalls)
	}

	if cache.deleteCalls != 0 {
		t.Fatalf("cache should not be invalidated, got %d calls", cache.deleteCalls)
	}
}

func TestSaveGroupMessage_RejectsEmptyNonce(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	service := newTestService(repo, cache)

	_, err := service.SaveGroupMessage(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
		"alice",
		"ciphertext",
		"",
		map[string]string{
			uuid.Must(uuid.NewV4()).String(): "key",
		},
	)

	if err == nil {
		t.Fatal("expected error for empty nonce")
	}

	if repo.saveCalls != 0 {
		t.Fatalf("repository should not be called")
	}
}

func TestSaveGroupMessage_RejectsEmptyKeys(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	service := newTestService(repo, cache)

	_, err := service.SaveGroupMessage(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
		"alice",
		"ciphertext",
		"nonce",
		map[string]string{},
	)

	if err == nil {
		t.Fatal("expected error for empty encrypted keys")
	}

	if repo.saveCalls != 0 {
		t.Fatalf("repository should not be called")
	}
}

func TestSaveGroupMessage_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database failure")

	repo := &mockMessageRepository{
		saveMessageFn: func(
			ctx context.Context,
			msg *domainMsg.Message,
			keys []domainMsg.MessageKey,
		) error {
			return expectedErr
		},
	}

	cache := newMockCache()
	service := newTestService(repo, cache)

	_, err := service.SaveGroupMessage(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
		"alice",
		"ciphertext",
		"nonce",
		map[string]string{
			uuid.Must(uuid.NewV4()).String(): "key",
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error %v, got %v", expectedErr, err)
	}

	if cache.deleteCalls != 0 {
		t.Fatalf("cache should not be invalidated when save fails")
	}
}

func TestSaveGroupMessage_InvalidRecipientUUIDIsIgnored(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	service := newTestService(repo, cache)

	validRecipient := uuid.Must(uuid.NewV4())

	_, err := service.SaveGroupMessage(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
		"alice",
		"ciphertext",
		"nonce",
		map[string]string{
			validRecipient.String(): "valid-key",
			"not-a-uuid":            "invalid-key",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.savedKeys) != 1 {
		t.Fatalf(
			"expected invalid UUID to be ignored and 1 valid key saved, got %d",
			len(repo.savedKeys),
		)
	}

	if repo.savedKeys[0].RecipientID != validRecipient {
		t.Fatalf("unexpected recipient ID %s", repo.savedKeys[0].RecipientID)
	}
}

func TestGetHistory_CacheMissLoadsDatabase(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	roomID := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())

	message1 := newTestMessage(
		uuid.Must(uuid.NewV4()),
		roomID,
		uuid.Must(uuid.NewV4()),
		"ciphertext-1",
		"nonce-1",
	)

	message2 := newTestMessage(
		uuid.Must(uuid.NewV4()),
		roomID,
		uuid.Must(uuid.NewV4()),
		"ciphertext-2",
		"nonce-2",
	)

	repo.getRoomHistoryFn = func(
		ctx context.Context,
		gotRoomID uuid.UUID,
		limit int,
	) ([]*domainMsg.Message, error) {
		if gotRoomID != roomID {
			t.Fatalf("expected room %s, got %s", roomID, gotRoomID)
		}

		if limit != 100 {
			t.Fatalf("expected history limit 100, got %d", limit)
		}

		return []*domainMsg.Message{message1, message2}, nil
	}

	repo.getMessageKeysForUserFn = func(
		ctx context.Context,
		messageIDs []uuid.UUID,
		gotUserID uuid.UUID,
	) (map[uuid.UUID]string, error) {
		if gotUserID != userID {
			t.Fatalf("expected user %s, got %s", userID, gotUserID)
		}

		return map[uuid.UUID]string{
			message1.ID: "key-1",
			message2.ID: "key-2",
		}, nil
	}

	service := newTestService(repo, cache)

	result, err := service.GetHistory(
		context.Background(),
		roomID,
		userID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(result))
	}

	if result[0].Message.ID != message1.ID {
		t.Fatalf("unexpected first message ID")
	}

	if result[0].EncryptedKey != "key-1" {
		t.Fatalf("expected key-1, got %s", result[0].EncryptedKey)
	}

	if result[1].Message.ID != message2.ID {
		t.Fatalf("unexpected second message ID")
	}

	if result[1].EncryptedKey != "key-2" {
		t.Fatalf("expected key-2, got %s", result[1].EncryptedKey)
	}

	if repo.historyCalls != 1 {
		t.Fatalf("expected 1 history query, got %d", repo.historyCalls)
	}

	if repo.keysCalls != 1 {
		t.Fatalf("expected 1 key query, got %d", repo.keysCalls)
	}

	if cache.getCalls != 1 {
		t.Fatalf("expected 1 cache read, got %d", cache.getCalls)
	}

	if cache.setCalls != 1 {
		t.Fatalf("expected 1 cache write, got %d", cache.setCalls)
	}
}

func TestGetHistory_CacheHitDoesNotLoadDatabaseHistory(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	roomID := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())

	message1 := newTestMessage(
		uuid.Must(uuid.NewV4()),
		roomID,
		uuid.Must(uuid.NewV4()),
		"ciphertext",
		"nonce",
	)

	cachedMessages := []*domainMsg.Message{
		message1,
	}

	data, err := json.Marshal(cachedMessages)
	if err != nil {
		t.Fatalf("failed to marshal cache: %v", err)
	}

	cache.values["history:"+roomID.String()] = data

	repo.getMessageKeysForUserFn = func(
		ctx context.Context,
		messageIDs []uuid.UUID,
		gotUserID uuid.UUID,
	) (map[uuid.UUID]string, error) {
		return map[uuid.UUID]string{
			message1.ID: "encrypted-key",
		}, nil
	}

	service := newTestService(repo, cache)

	result, err := service.GetHistory(
		context.Background(),
		roomID,
		userID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 message, got %d", len(result))
	}

	if result[0].EncryptedKey != "encrypted-key" {
		t.Fatalf("unexpected encrypted key")
	}

	if repo.historyCalls != 0 {
		t.Fatalf(
			"database history should not be queried on cache hit, got %d calls",
			repo.historyCalls,
		)
	}

	if cache.getCalls != 1 {
		t.Fatalf("expected 1 cache read, got %d", cache.getCalls)
	}
}

func TestGetHistory_OnlyReturnsMessagesWithUserKey(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	roomID := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())

	message1 := newTestMessage(
		uuid.Must(uuid.NewV4()),
		roomID,
		uuid.Must(uuid.NewV4()),
		"ciphertext-1",
		"nonce-1",
	)

	message2 := newTestMessage(
		uuid.Must(uuid.NewV4()),
		roomID,
		uuid.Must(uuid.NewV4()),
		"ciphertext-2",
		"nonce-2",
	)

	repo.getRoomHistoryFn = func(
		ctx context.Context,
		roomID uuid.UUID,
		limit int,
	) ([]*domainMsg.Message, error) {
		return []*domainMsg.Message{
			message1,
			message2,
		}, nil
	}

	repo.getMessageKeysForUserFn = func(
		ctx context.Context,
		messageIDs []uuid.UUID,
		userID uuid.UUID,
	) (map[uuid.UUID]string, error) {
		return map[uuid.UUID]string{
			message1.ID: "key-for-user",
		}, nil
	}

	service := newTestService(repo, cache)

	result, err := service.GetHistory(
		context.Background(),
		roomID,
		userID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf(
			"expected only messages with recipient key, got %d",
			len(result),
		)
	}

	if result[0].Message.ID != message1.ID {
		t.Fatalf("unexpected message returned")
	}
}

func TestGetHistory_HistoryRepositoryError(t *testing.T) {
	expectedErr := errors.New("history database failure")

	repo := &mockMessageRepository{
		getRoomHistoryFn: func(
			ctx context.Context,
			roomID uuid.UUID,
			limit int,
		) ([]*domainMsg.Message, error) {
			return nil, expectedErr
		},
	}

	cache := newMockCache()
	service := newTestService(repo, cache)

	_, err := service.GetHistory(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestGetHistory_KeyRepositoryError(t *testing.T) {
	expectedErr := errors.New("key lookup failure")

	roomID := uuid.Must(uuid.NewV4())

	repo := &mockMessageRepository{
		getRoomHistoryFn: func(
			ctx context.Context,
			roomID uuid.UUID,
			limit int,
		) ([]*domainMsg.Message, error) {
			return []*domainMsg.Message{
				newTestMessage(
					uuid.Must(uuid.NewV4()),
					roomID,
					uuid.Must(uuid.NewV4()),
					"ciphertext",
					"nonce",
				),
			}, nil
		},

		getMessageKeysForUserFn: func(
			ctx context.Context,
			messageIDs []uuid.UUID,
			userID uuid.UUID,
		) (map[uuid.UUID]string, error) {
			return nil, expectedErr
		},
	}

	cache := newMockCache()
	service := newTestService(repo, cache)

	_, err := service.GetHistory(
		context.Background(),
		roomID,
		uuid.Must(uuid.NewV4()),
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestGetHistory_UsesCorrectRoomAndUser(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	roomID := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())

	messageID := uuid.Must(uuid.NewV4())

	repo.getRoomHistoryFn = func(
		ctx context.Context,
		gotRoomID uuid.UUID,
		limit int,
	) ([]*domainMsg.Message, error) {
		if gotRoomID != roomID {
			t.Fatalf("wrong room ID: expected %s, got %s", roomID, gotRoomID)
		}

		return []*domainMsg.Message{
			newTestMessage(
				messageID,
				roomID,
				uuid.Must(uuid.NewV4()),
				"ciphertext",
				"nonce",
			),
		}, nil
	}

	repo.getMessageKeysForUserFn = func(
		ctx context.Context,
		messageIDs []uuid.UUID,
		gotUserID uuid.UUID,
	) (map[uuid.UUID]string, error) {
		if gotUserID != userID {
			t.Fatalf("wrong user ID: expected %s, got %s", userID, gotUserID)
		}

		if len(messageIDs) != 1 || messageIDs[0] != messageID {
			t.Fatalf("unexpected message IDs")
		}

		return map[uuid.UUID]string{
			messageID: "key",
		}, nil
	}

	service := newTestService(repo, cache)

	_, err := service.GetHistory(
		context.Background(),
		roomID,
		userID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteMessage_Success(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	messageID := uuid.Must(uuid.NewV4())
	roomID := uuid.Must(uuid.NewV4())
	senderID := uuid.Must(uuid.NewV4())

	msg := newTestMessage(
		messageID,
		roomID,
		senderID,
		"ciphertext",
		"nonce",
	)

	repo.getMessageByIDFn = func(
		ctx context.Context,
		id uuid.UUID,
	) (*domainMsg.Message, error) {
		if id != messageID {
			t.Fatalf("expected message ID %s, got %s", messageID, id)
		}

		return msg, nil
	}

	service := newTestService(repo, cache)

	err := service.DeleteMessage(
		context.Background(),
		messageID,
		senderID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.getMessageCalls != 1 {
		t.Fatalf("expected 1 GetMessageByID call")
	}

	if repo.deleteCalls != 1 {
		t.Fatalf("expected 1 DeleteMessage call")
	}

	if repo.lastDeletedMessageID != messageID {
		t.Fatalf("wrong deleted message ID")
	}

	if cache.deleteCalls != 1 {
		t.Fatalf("expected cache invalidation")
	}

	expectedKey := "history:" + roomID.String()

	if cache.lastDeleteKey != expectedKey {
		t.Fatalf(
			"expected cache key %q, got %q",
			expectedKey,
			cache.lastDeleteKey,
		)
	}
}

func TestDeleteMessage_ForbiddenForDifferentUser(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	messageID := uuid.Must(uuid.NewV4())
	roomID := uuid.Must(uuid.NewV4())
	senderID := uuid.Must(uuid.NewV4())
	attackerID := uuid.Must(uuid.NewV4())

	msg := newTestMessage(
		messageID,
		roomID,
		senderID,
		"ciphertext",
		"nonce",
	)

	repo.getMessageByIDFn = func(
		ctx context.Context,
		id uuid.UUID,
	) (*domainMsg.Message, error) {
		return msg, nil
	}

	service := newTestService(repo, cache)

	err := service.DeleteMessage(
		context.Background(),
		messageID,
		attackerID,
	)

	if err == nil {
		t.Fatal("expected forbidden error")
	}

	if repo.deleteCalls != 0 {
		t.Fatalf(
			"message must not be deleted by another user, got %d calls",
			repo.deleteCalls,
		)
	}

	if cache.deleteCalls != 0 {
		t.Fatalf(
			"cache should not be invalidated when delete is forbidden",
		)
	}
}

func TestDeleteMessage_GetMessageError(t *testing.T) {
	expectedErr := errors.New("message not found")

	repo := &mockMessageRepository{
		getMessageByIDFn: func(
			ctx context.Context,
			id uuid.UUID,
		) (*domainMsg.Message, error) {
			return nil, expectedErr
		},
	}

	cache := newMockCache()
	service := newTestService(repo, cache)

	err := service.DeleteMessage(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if repo.deleteCalls != 0 {
		t.Fatalf("DeleteMessage should not be called")
	}
}

func TestDeleteMessage_RepositoryError(t *testing.T) {
	expectedErr := errors.New("delete database failure")

	messageID := uuid.Must(uuid.NewV4())
	roomID := uuid.Must(uuid.NewV4())
	senderID := uuid.Must(uuid.NewV4())

	repo := &mockMessageRepository{
		getMessageByIDFn: func(
			ctx context.Context,
			id uuid.UUID,
		) (*domainMsg.Message, error) {
			return newTestMessage(
				messageID,
				roomID,
				senderID,
				"ciphertext",
				"nonce",
			), nil
		},

		deleteMessageFn: func(
			ctx context.Context,
			id uuid.UUID,
		) error {
			return expectedErr
		},
	}

	cache := newMockCache()
	service := newTestService(repo, cache)

	err := service.DeleteMessage(
		context.Background(),
		messageID,
		senderID,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if cache.deleteCalls != 0 {
		t.Fatalf(
			"cache should not be invalidated when repository delete fails",
		)
	}
}

func TestEditMessage_Success(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	messageID := uuid.Must(uuid.NewV4())
	roomID := uuid.Must(uuid.NewV4())
	senderID := uuid.Must(uuid.NewV4())
	recipientID := uuid.Must(uuid.NewV4())

	msg := newTestMessage(
		messageID,
		roomID,
		senderID,
		"old-ciphertext",
		"old-nonce",
	)

	repo.getMessageByIDFn = func(
		ctx context.Context,
		id uuid.UUID,
	) (*domainMsg.Message, error) {
		return msg, nil
	}

	service := newTestService(repo, cache)

	err := service.EditMessage(
		context.Background(),
		messageID,
		senderID,
		"new-ciphertext",
		"new-nonce",
		map[string]string{
			recipientID.String(): "new-encrypted-key",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.getMessageCalls != 1 {
		t.Fatalf("expected 1 GetMessageByID call")
	}

	if repo.editCalls != 1 {
		t.Fatalf("expected 1 EditMessage call")
	}

	if repo.lastEditedMessageID != messageID {
		t.Fatalf("wrong message ID")
	}

	if repo.lastCiphertext != "new-ciphertext" {
		t.Fatalf(
			"expected new ciphertext, got %s",
			repo.lastCiphertext,
		)
	}

	if repo.lastNonce != "new-nonce" {
		t.Fatalf(
			"expected new nonce, got %s",
			repo.lastNonce,
		)
	}

	if len(repo.lastEditedKeys) != 1 {
		t.Fatalf("expected 1 replacement key")
	}

	if repo.lastEditedKeys[0].MessageID != messageID {
		t.Fatalf("wrong MessageID in replacement key")
	}

	if repo.lastEditedKeys[0].RecipientID != recipientID {
		t.Fatalf("wrong RecipientID in replacement key")
	}

	if repo.lastEditedKeys[0].EncryptedKey != "new-encrypted-key" {
		t.Fatalf("wrong encrypted key")
	}

	if cache.deleteCalls != 1 {
		t.Fatalf("expected cache invalidation")
	}

	expectedKey := "history:" + roomID.String()

	if cache.lastDeleteKey != expectedKey {
		t.Fatalf(
			"expected cache key %q, got %q",
			expectedKey,
			cache.lastDeleteKey,
		)
	}
}

func TestEditMessage_RejectsEmptyCiphertext(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	service := newTestService(repo, cache)

	err := service.EditMessage(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
		"",
		"nonce",
		map[string]string{
			uuid.Must(uuid.NewV4()).String(): "key",
		},
	)

	if err == nil {
		t.Fatal("expected error for empty ciphertext")
	}

	if repo.getMessageCalls != 0 {
		t.Fatalf("repository should not be called")
	}
}

func TestEditMessage_RejectsEmptyNonce(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	service := newTestService(repo, cache)

	err := service.EditMessage(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
		"ciphertext",
		"",
		map[string]string{
			uuid.Must(uuid.NewV4()).String(): "key",
		},
	)

	if err == nil {
		t.Fatal("expected error for empty nonce")
	}

	if repo.getMessageCalls != 0 {
		t.Fatalf("repository should not be called")
	}
}

func TestEditMessage_RejectsEmptyKeys(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	service := newTestService(repo, cache)

	err := service.EditMessage(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
		"ciphertext",
		"nonce",
		map[string]string{},
	)

	if err == nil {
		t.Fatal("expected error for empty encrypted keys")
	}

	if repo.getMessageCalls != 0 {
		t.Fatalf("repository should not be called")
	}
}

func TestEditMessage_ForbiddenForDifferentUser(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	messageID := uuid.Must(uuid.NewV4())
	roomID := uuid.Must(uuid.NewV4())
	senderID := uuid.Must(uuid.NewV4())
	attackerID := uuid.Must(uuid.NewV4())

	msg := newTestMessage(
		messageID,
		roomID,
		senderID,
		"old-ciphertext",
		"old-nonce",
	)

	repo.getMessageByIDFn = func(
		ctx context.Context,
		id uuid.UUID,
	) (*domainMsg.Message, error) {
		return msg, nil
	}

	service := newTestService(repo, cache)

	err := service.EditMessage(
		context.Background(),
		messageID,
		attackerID,
		"new-ciphertext",
		"new-nonce",
		map[string]string{
			uuid.Must(uuid.NewV4()).String(): "new-key",
		},
	)

	if err == nil {
		t.Fatal("expected forbidden error")
	}

	if repo.editCalls != 0 {
		t.Fatalf(
			"message must not be edited by another user, got %d calls",
			repo.editCalls,
		)
	}

	if cache.deleteCalls != 0 {
		t.Fatalf(
			"cache should not be invalidated when edit is forbidden",
		)
	}
}

func TestEditMessage_GetMessageError(t *testing.T) {
	expectedErr := errors.New("message not found")

	repo := &mockMessageRepository{
		getMessageByIDFn: func(
			ctx context.Context,
			id uuid.UUID,
		) (*domainMsg.Message, error) {
			return nil, expectedErr
		},
	}

	cache := newMockCache()
	service := newTestService(repo, cache)

	err := service.EditMessage(
		context.Background(),
		uuid.Must(uuid.NewV4()),
		uuid.Must(uuid.NewV4()),
		"ciphertext",
		"nonce",
		map[string]string{
			uuid.Must(uuid.NewV4()).String(): "key",
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if repo.editCalls != 0 {
		t.Fatalf("EditMessage repository method should not be called")
	}
}

func TestEditMessage_RepositoryError(t *testing.T) {
	expectedErr := errors.New("edit database failure")

	messageID := uuid.Must(uuid.NewV4())
	roomID := uuid.Must(uuid.NewV4())
	senderID := uuid.Must(uuid.NewV4())

	repo := &mockMessageRepository{
		getMessageByIDFn: func(
			ctx context.Context,
			id uuid.UUID,
		) (*domainMsg.Message, error) {
			return newTestMessage(
				messageID,
				roomID,
				senderID,
				"old-ciphertext",
				"old-nonce",
			), nil
		},

		editMessageFn: func(
			ctx context.Context,
			id uuid.UUID,
			ciphertext string,
			nonce string,
			keys []domainMsg.MessageKey,
		) error {
			return expectedErr
		},
	}

	cache := newMockCache()
	service := newTestService(repo, cache)

	err := service.EditMessage(
		context.Background(),
		messageID,
		senderID,
		"new-ciphertext",
		"new-nonce",
		map[string]string{
			uuid.Must(uuid.NewV4()).String(): "new-key",
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}

	if cache.deleteCalls != 0 {
		t.Fatalf(
			"cache should not be invalidated when edit fails",
		)
	}
}

func TestEditMessage_InvalidRecipientUUIDIsIgnored(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	messageID := uuid.Must(uuid.NewV4())
	roomID := uuid.Must(uuid.NewV4())
	senderID := uuid.Must(uuid.NewV4())
	validRecipient := uuid.Must(uuid.NewV4())

	repo.getMessageByIDFn = func(
		ctx context.Context,
		id uuid.UUID,
	) (*domainMsg.Message, error) {
		return newTestMessage(
			messageID,
			roomID,
			senderID,
			"old",
			"old-nonce",
		), nil
	}

	service := newTestService(repo, cache)

	err := service.EditMessage(
		context.Background(),
		messageID,
		senderID,
		"new",
		"new-nonce",
		map[string]string{
			validRecipient.String(): "valid-key",
			"not-a-valid-uuid":      "invalid-key",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.lastEditedKeys) != 1 {
		t.Fatalf(
			"expected invalid UUID to be ignored and one valid key, got %d",
			len(repo.lastEditedKeys),
		)
	}

	if repo.lastEditedKeys[0].RecipientID != validRecipient {
		t.Fatalf("unexpected recipient ID")
	}
}

func TestDeleteMessage_DoesNotDeleteDifferentRoomCache(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	messageID := uuid.Must(uuid.NewV4())
	roomID := uuid.Must(uuid.NewV4())
	senderID := uuid.Must(uuid.NewV4())

	repo.getMessageByIDFn = func(
		ctx context.Context,
		id uuid.UUID,
	) (*domainMsg.Message, error) {
		return newTestMessage(
			messageID,
			roomID,
			senderID,
			"ciphertext",
			"nonce",
		), nil
	}

	service := newTestService(repo, cache)

	err := service.DeleteMessage(
		context.Background(),
		messageID,
		senderID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedKey := "history:" + roomID.String()

	if cache.lastDeleteKey != expectedKey {
		t.Fatalf(
			"expected only message room cache to be invalidated: %q, got %q",
			expectedKey,
			cache.lastDeleteKey,
		)
	}
}

func TestHistoryCacheKeyUsesRoomID(t *testing.T) {
	repo := &mockMessageRepository{}
	cache := newMockCache()

	roomID := uuid.Must(uuid.NewV4())

	repo.getRoomHistoryFn = func(
		ctx context.Context,
		roomID uuid.UUID,
		limit int,
	) ([]*domainMsg.Message, error) {
		return []*domainMsg.Message{}, nil
	}

	service := newTestService(repo, cache)

	_, err := service.GetHistory(
		context.Background(),
		roomID,
		uuid.Must(uuid.NewV4()),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedKey := "history:" + roomID.String()

	if cache.lastGetKey != expectedKey {
		t.Fatalf(
			"expected cache key %q, got %q",
			expectedKey,
			cache.lastGetKey,
		)
	}

	if cache.lastSetKey != expectedKey {
		t.Fatalf(
			"expected cache set key %q, got %q",
			expectedKey,
			cache.lastSetKey,
		)
	}
}

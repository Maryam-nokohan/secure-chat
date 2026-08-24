package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type UserSummary struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Online    bool      `json:"online"`
}

type Stats struct {
	TotalUsers    int64 `json:"total_users"`
	OnlineUsers   int   `json:"online_users"`
	TotalRooms    int64 `json:"total_rooms"`
	TotalMessages int64 `json:"total_messages"`
	UptimeSeconds int64 `json:"uptime_seconds"`
}

type Service struct {
	userRepo  ports.UserRepository
	chatRepo  ports.ChatRepositoryI
	msgRepo   ports.MessageRepository
	invoker   *Invoker
	onlineFn  func() int
	startedAt time.Time
}

func NewService(userRepo ports.UserRepository, chatRepo ports.ChatRepositoryI, msgRepo ports.MessageRepository, onlineFn func() int) *Service {
	return &Service{
		userRepo:  userRepo,
		chatRepo:  chatRepo,
		msgRepo:   msgRepo,
		invoker:   NewInvoker(100),
		onlineFn:  onlineFn,
		startedAt: time.Now(),
	}
}


func (s *Service) ListUsers(ctx context.Context, limit, offset int) ([]UserSummary, error) {
	users, err := s.userRepo.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]UserSummary, len(users))
	for i, u := range users {
		out[i] = UserSummary{ID: u.ID.String(), Username: u.Username, Role: u.Role, CreatedAt: u.CreatedAt}
	}
	return out, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID, actor string) error {
	if err := s.invoker.Do(ctx, &DeleteUserCommand{Repo: s.userRepo, UserID: id}); err != nil {
		return err
	}
	pkg.LogInfo(fmt.Sprintf("[admin:%s] deleted user %s", actor, id))
	return nil
}

func (s *Service) EditUser(ctx context.Context, id uuid.UUID, newUsername, newBio, actor string) error {
	cmd := &EditUserCommand{Repo: s.userRepo, UserID: id, NewUsername: newUsername, NewBio: newBio}
	if err := s.invoker.Do(ctx, cmd); err != nil {
		return err
	}
	pkg.LogInfo(fmt.Sprintf("[admin:%s] edited user %s", actor, id))
	return nil
}

func (s *Service) SetRole(ctx context.Context, id uuid.UUID, role, actor string) error {
	if role != "admin" && role != "user" {
		return errors.New("role must be 'admin' or 'user'")
	}
	cmd := &SetRoleCommand{Repo: s.userRepo, UserID: id, NewRole: role}
	if err := s.invoker.Do(ctx, cmd); err != nil {
		return err
	}
	pkg.LogInfo(fmt.Sprintf("[admin:%s] set role of %s to %s", actor, id, role))
	return nil
}

func (s *Service) CreateAdmin(ctx context.Context, username, passHash, publicKey string) (*user.User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	u := user.User{ID: id, Username: username, PassHash: passHash, PublicKey: publicKey, Role: "admin"}
	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) Undo(ctx context.Context) (string, error) { return s.invoker.UndoLast(ctx) }
func (s *Service) History() []HistoryEntry                  { return s.invoker.History() }

func (s *Service) Stats(ctx context.Context) (Stats, error) {
	totalUsers, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		return Stats{}, err
	}
	totalRooms, err := s.chatRepo.CountRooms(ctx)
	if err != nil {
		return Stats{}, err
	}
	totalMessages, err := s.msgRepo.CountMessages(ctx)
	if err != nil {
		return Stats{}, err
	}
	online := 0
	if s.onlineFn != nil {
		online = s.onlineFn()
	}
	return Stats{
		TotalUsers: totalUsers, OnlineUsers: online, TotalRooms: totalRooms,
		TotalMessages: totalMessages, UptimeSeconds: int64(time.Since(s.startedAt).Seconds()),
	}, nil
}
package contact

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	domainContact "github.com/maryam-nokohan/secure-chat/internal/core/domain/contact"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type Service struct {
	repo     ports.ContactRepository
	userRepo ports.UserRepository
	chatSvc  ports.ChatServiceI
}

func NewService(repo ports.ContactRepository, userRepo ports.UserRepository, chatSvc ports.ChatServiceI) ports.ContactServiceI {
	pkg.LogInfo("Init ContactService...")
	return &Service{repo: repo, userRepo: userRepo, chatSvc: chatSvc}
}

func (s *Service) SendRequest(ctx context.Context, requesterID uuid.UUID, targetPublicID string) (*domainContact.Contact, error) {
	targetPublicID = strings.ToUpper(strings.TrimSpace(targetPublicID))
	if targetPublicID == "" {
		return nil, errors.New("contact id is required")
	}

	target, err := s.userRepo.FindUserByPublicID(ctx, targetPublicID)
	if err != nil {
		return nil, errors.New("no user found with that contact id")
	}
	if target.ID == requesterID {
		return nil, errors.New("you cannot add yourself")
	}

	existing, err := s.repo.FindByPair(ctx, requesterID, target.ID)
	if err == nil && existing != nil {
		switch existing.Status {
		case domainContact.StatusAccepted:
			return nil, errors.New("already contacts")
		case domainContact.StatusBlocked:
			return nil, errors.New("unable to send request")
		case domainContact.StatusPending:
			if existing.RequesterID == requesterID {
				return nil, errors.New("request already sent")
			}
			return s.Accept(ctx, existing.ID, requesterID)
		case domainContact.StatusDeclined:
			if err := s.repo.UpdateStatus(ctx, existing.ID, domainContact.StatusPending, nil); err != nil {
				return nil, err
			}
			existing.Status = domainContact.StatusPending
			existing.RequesterID = requesterID
			existing.AddresseeID = target.ID
			return existing, nil
		}
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	c := &domainContact.Contact{
		ID: id, RequesterID: requesterID, AddresseeID: target.ID,
		Status: domainContact.StatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) Accept(ctx context.Context, contactID, userID uuid.UUID) (*domainContact.Contact, error) {
	c, err := s.repo.FindByID(ctx, contactID)
	if err != nil {
		return nil, errors.New("request not found")
	}
	if c.AddresseeID != userID {
		return nil, errors.New("only the recipient can accept this request")
	}
	if c.Status != domainContact.StatusPending {
		return nil, errors.New("request is not pending")
	}

	room, err := s.chatSvc.CreateDirectRoom(ctx, c.RequesterID, c.AddresseeID)
	if err != nil {
		pkg.LogError(err)
		return nil, errors.New("failed to set up the conversation")
	}
	if err := s.repo.UpdateStatus(ctx, c.ID, domainContact.StatusAccepted, &room.ID); err != nil {
		return nil, err
	}
	c.Status = domainContact.StatusAccepted
	c.RoomID = &room.ID
	return c, nil
}

func (s *Service) Decline(ctx context.Context, contactID, userID uuid.UUID) error {
	c, err := s.repo.FindByID(ctx, contactID)
	if err != nil {
		return errors.New("request not found")
	}
	if c.AddresseeID != userID && c.RequesterID != userID {
		return errors.New("not authorized")
	}
	if c.Status != domainContact.StatusPending {
		return errors.New("request is not pending")
	}
	return s.repo.UpdateStatus(ctx, c.ID, domainContact.StatusDeclined, nil)
}

func (s *Service) Block(ctx context.Context, contactID, userID uuid.UUID) error {
	c, err := s.repo.FindByID(ctx, contactID)
	if err != nil {
		return errors.New("request not found")
	}
	if c.AddresseeID != userID && c.RequesterID != userID {
		return errors.New("not authorized")
	}
	return s.repo.UpdateStatus(ctx, c.ID, domainContact.StatusBlocked, nil)
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]ports.ContactSummary, error) {
	rows, err := s.repo.ListForUser(ctx, userID, "")
	if err != nil {
		return nil, err
	}
	out := make([]ports.ContactSummary, 0, len(rows))
	for _, c := range rows {
		otherID := c.RequesterID
		incoming := true
		if c.RequesterID == userID {
			otherID = c.AddresseeID
			incoming = false
		}
		other, err := s.userRepo.FindUserByID(ctx, otherID)
		if err != nil {
			continue
		}
		roomID := ""
		if c.RoomID != nil {
			roomID = c.RoomID.String()
		}
		out = append(out, ports.ContactSummary{
			ID: c.ID.String(), User: *other, Status: c.Status, Incoming: incoming, RoomID: roomID,
		})
	}
	return out, nil
}

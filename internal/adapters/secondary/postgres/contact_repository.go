package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/contact"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"gorm.io/gorm"
)

type ContactRepo struct{ db *gorm.DB }

func NewContactRepository(db *gorm.DB) ports.ContactRepository {
	pkg.LogInfo("Initializing ContactRepository...")
	return &ContactRepo{db: db}
}

func (r *ContactRepo) Create(ctx context.Context, c *contact.Contact) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *ContactRepo) FindByID(ctx context.Context, id uuid.UUID) (*contact.Contact, error) {
	var c contact.Contact
	if err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContactRepo) FindByPair(ctx context.Context, a, b uuid.UUID) (*contact.Contact, error) {
	var c contact.Contact
	err := r.db.WithContext(ctx).
		Where("(requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)", a, b, b, a).
		First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContactRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status contact.Status, roomID *uuid.UUID) error {
	updates := map[string]interface{}{"status": status, "updated_at": gorm.Expr("now()")}
	if roomID != nil {
		updates["room_id"] = *roomID
	}
	return r.db.WithContext(ctx).Model(&contact.Contact{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ContactRepo) ListForUser(ctx context.Context, userID uuid.UUID, status contact.Status) ([]*contact.Contact, error) {
	q := r.db.WithContext(ctx).Where("requester_id = ? OR addressee_id = ?", userID, userID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var list []*contact.Contact
	if err := q.Order("updated_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
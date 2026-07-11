package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"gorm.io/gorm"
)

const userCacheTTL = 10 * time.Minute

type UserRepository struct {
	db    *gorm.DB
	cache ports.Cache
}

func NewUserRepositoryService(db *gorm.DB, cache ports.Cache) (ports.UserRepository, error) {
	pkg.LogInfo("Initializing UserRepository...")
	return &UserRepository{db: db, cache: cache}, nil
}

func userIDKey(id uuid.UUID) string    { return "user:id:" + id.String() }
func userNameKey(name string) string   { return "user:name:" + name }

func (r *UserRepository) CreateUser(ctx context.Context, u user.User) error {
	pkg.LogRepo("Creating User " + u.Username)
	if err := r.db.WithContext(ctx).Create(&u).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	key := userIDKey(id)
	if cached, err := r.cache.Get(ctx, key); err == nil {
		var u user.User
		if json.Unmarshal(cached, &u) == nil {
			return &u, nil
		}
	}

	var dbModel user.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&dbModel).Error; err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, dbModel, userCacheTTL)
	return &dbModel, nil
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
	key := userNameKey(username)
	if cached, err := r.cache.Get(ctx, key); err == nil {
		var u user.User
		if json.Unmarshal(cached, &u) == nil {
			return &u, nil
		}
	}

	var u user.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	_ = r.cache.Set(ctx, key, u, userCacheTTL)
	return &u, nil
}

func (r *UserRepository) EditUser(ctx context.Context, u user.User) error {
	pkg.LogRepo("Editing User " + u.Username)
	if err := r.db.WithContext(ctx).Save(&u).Error; err != nil {
		return err
	}
	_ = r.cache.Delete(ctx, userIDKey(u.ID))
	_ = r.cache.Delete(ctx, userNameKey(u.Username))
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, u user.User) error {
	pkg.LogRepo("Deleting user " + u.Username)
	if err := r.db.WithContext(ctx).Delete(&u).Error; err != nil {
		return err
	}
	_ = r.cache.Delete(ctx, userIDKey(u.ID))
	_ = r.cache.Delete(ctx, userNameKey(u.Username))
	return nil
}
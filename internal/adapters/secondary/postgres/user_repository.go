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

type cachedUser struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	PublicKey string    `json:"public_key"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toCached(u user.User) cachedUser {
	return cachedUser{ID: u.ID, Username: u.Username, PublicKey: u.PublicKey, Bio: u.Bio, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt}
}
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
	
	var u user.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
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
func (r *UserRepository) ListUsers(ctx context.Context, limit, offset int) ([]user.User, error) {
	var users []user.User
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.User{}).Count(&count).Error
	return count, err
}

func (r *UserRepository) RestoreUser(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Unscoped().Model(&user.User{}).
		Where("id = ?", id).Update("deleted_at", nil).Error
	if err != nil {
		return err
	}
	_ = r.cache.Delete(ctx, userIDKey(id))
	return nil
}
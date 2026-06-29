package postgres

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepositoryService(db *gorm.DB) (ports.UserRepository, error) {
	pkg.LogInfo("Initializing UserRepository...")
	return &UserRepository{
		db: db,
	}, nil

}

func (r *UserRepository) CreateUser(ctx context.Context, user user.User) error {
	pkg.LogRepo("Creating User " + user.Username)
	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *UserRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var dbModel user.User
	
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&dbModel)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return &dbModel, nil
}
func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	
	err := r.db.
	WithContext(ctx).
	Where("username = ?", username).
	First(&u).
	Error
	
	if err != nil {
		return nil, err
	}
	
	return &u, nil
}
func (r *UserRepository) EditUser(ctx context.Context, user user.User) error {
	pkg.LogRepo("Editing User " + user.Username)
	
	return r.db.
	WithContext(ctx).
		Save(&user).
		Error
	}
func (r *UserRepository) DeleteUser(ctx context.Context, user user.User) error {
	pkg.LogRepo("Deleting user " + user.Username)
	return r.db.
		WithContext(ctx).
		Delete(&user).
		Error
}

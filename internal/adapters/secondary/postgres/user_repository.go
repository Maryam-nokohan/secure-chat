package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepositoryService(db *gorm.DB) (ports.UserRepository, error) {

	return &UserRepository{
		db: db,
	}, nil

}

func (r *UserRepository) CreateUser(ctx context.Context, user user.User) error {

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

	return &user.User{
		ID:       dbModel.ID,
		Username: dbModel.Username,
		PassHash: dbModel.PassHash,
		Bio: dbModel.Bio,
	}, nil
}
func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {

	var dbModel user.User

	result := r.db.WithContext(ctx).Where("username = ?", username).First(&dbModel)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user.User{
		ID:       dbModel.ID,
		Username: dbModel.Username,
		PassHash: dbModel.PassHash,
		Bio: dbModel.Bio,
	}, nil
}
func (r *UserRepository) EditUser(ctx context.Context, user user.User) error {

	if err := r.db.WithContext(ctx).Where("id = ?", user.ID).Error; err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return err
	}
	return nil

}
func (r *UserRepository) DeleteUser(ctx context.Context, user user.User) error {

	result := r.db.WithContext(ctx).Where("id = ?", user.ID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found or couldn't be deleted")
	}
	return nil
}

package postgres

import (
	"gorm.io/driver/postgres"

	"github.com/maryam-nokohan/secure-chat/internal/domain/user"
	"github.com/maryam-nokohan/secure-chat/internal/ports/outbound"
	"gorm.io/gorm"
)

type UserRepositoryService struct {
	db *gorm.db
}

func NewUserRepositoryService(dsn string) (outbound.UserRepository, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &UserRepositoryService{
		db: db,
	}, nil

}

func (r *UserRepositoryService) CreateUser(u user.User) error {
	return r.db.Create(&u).Error
}

func (r *UserRepositoryService) FindUser(id string) (*user.User, error) {
	var u user.User

	err := r.db.First(&u, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryService) EditUser(u user.User) error {
	return r.db.Save(&u).Error
}

func (r *UserRepositoryService) DeleteUser(u user.User) error {
	return r.db.Delete(&u).Error
}

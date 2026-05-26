package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUSerService(repo ports.UserRepository) ports.UserServicesI {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(name, password string) error {
	if exist, ok := s.repo.FindUserByUsername(context.Background(), name); ok == nil && exist != nil {
		return errors.New("this username already exists")
	}
	hash, err := pkg.HashPassword(password)
	if err != nil {
		return err
	}
	uID , err := uuid.NewV4()
	if err != nil {
		return  fmt.Errorf("faild to generate uuid %w" , err)
	}
	newUser := user.User{
		ID:          uID,
		Username:    name,
		PassHash:    hash,
		Bio:         "",
	}
	return s.repo.CreateUser(context.Background() , newUser)
}

func (s *UserService) Login(username, password string) error {
	u, err := s.repo.FindUserByUsername(context.Background(), username)
	if err != nil {
		return errors.New("Invalid credentials")
	}

	err = pkg.CheckPassword(password, u.PassHash)
	if err != nil {
		return errors.New("Invalid credentials")
	}
	fmt.Println("Login successfully")
	return nil
}


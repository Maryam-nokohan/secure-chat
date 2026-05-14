package user

import "errors"

type UserService struct {
	repo UserRepository
}

func NewUSerService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(name , passHash string) error {
	if _ , ok := s.repo.FindUser(name); ok {
		return errors.New("this username already exists")
	}
	newUser := User{
		Userame: name,
		PassHash: passHash,
		HasUsername: true,
		Bio: "",
		IsActive: true,
	}
	return s.repo.CreateUser(newUser)
}

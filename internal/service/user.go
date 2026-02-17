package service

import (
	"authapp/internal/storage/models"
	"authapp/internal/storage/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repository.UserRepo
}

func NewUserService(userRepo *repository.UserRepo) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) CreateUser(username, password string) (*models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(passwordHash),
	}

	err = s.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Authenticate(username, password string) (*models.User, error) {
	user, err := s.UserRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	print(user.PasswordHash, password)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

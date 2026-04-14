package services

import (
	"fmt"
	"phonebook_gorm/db"
	"phonebook_gorm/users"
)

type UserService struct {
	repo *users.Service
}

func NewUserService(repo *users.Service) *UserService {
	return &UserService{repo: repo}
}

// връща users
func (s *UserService) GetUsers() ([]db.User, error) {
	return s.repo.GetUsers()
}

// създава user + validation
func (s *UserService) CreateUser(user *db.User) error {

	if user.Username == "" {
		return fmt.Errorf("username is required")
	}

	if user.Email == "" {
		return fmt.Errorf("email is required")
	}

	return s.repo.CreateUser(user)
}

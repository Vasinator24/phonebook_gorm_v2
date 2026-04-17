package services

import (
	"phonebook_gorm/db"
	"phonebook_gorm/repository"
)

type UserService struct {
	repo *repository.Service
}

func NewUserService(repo *repository.Service) *UserService {
	return &UserService{repo: repo}
}

// GetUsers връща всички users
func (s *UserService) GetUsers() ([]db.User, error) {
	return s.repo.GetUsers()
}

// CreateUser създава нов user
func (s *UserService) CreateUser(user *db.User) error {
	return s.repo.CreateUser(user)
}

// UpdateUser обновява user
func (s *UserService) UpdateUser(user *db.User) error {
	return s.repo.UpdateUser(user)
}

// DeleteUser изтрива user по ID
func (s *UserService) DeleteUser(userID uint) error {
	return s.repo.DeleteUser(userID)
}

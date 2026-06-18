package services

import (
	"phonebook_gorm/db"
	"phonebook_gorm/repository"
	"time"
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

// UpdateUser обновява само позволените user полета
func (s *UserService) UpdateUser(id string, names string, email string) error {
	return s.repo.UpdateUser(id, names, email)
}

// DeleteUser изтрива user по ID
func (s *UserService) DeleteUser(userID uint) error {
	return s.repo.DeleteUser(userID)
}

func (s *UserService) GetUserByEmail(email string) (*db.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) GetUserByID(id uint) (*db.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) BlacklistToken(tokenHash string, expiresAt time.Time) error {
	return s.repo.BlacklistToken(tokenHash, expiresAt)
}

func (s *UserService) DeleteExpiredBlacklistedTokens() error {
	return s.repo.DeleteExpiredBlacklistedTokens()
}

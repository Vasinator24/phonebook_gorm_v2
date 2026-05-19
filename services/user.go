package services

import (
	"phonebook_gorm/db"
	"phonebook_gorm/repository"

	"golang.org/x/crypto/bcrypt"
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

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	// replace plain password with hash
	user.Password = string(hashedPassword)

	// save user
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

func (s *UserService) GetUserByEmail(email string) (*db.User, error) {
	return s.repo.GetUserByEmail(email)
}

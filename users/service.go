package users

import (
	"phonebook_gorm/db"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(dbConn *gorm.DB) *Service {
	return &Service{DB: dbConn}
}

// взима всички users
func (s *Service) GetUsers() ([]db.User, error) {
	var users []db.User

	err := s.DB.Preload("Phones").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// създава user
func (s *Service) CreateUser(user *db.User) error {
	return s.DB.Create(user).Error
}

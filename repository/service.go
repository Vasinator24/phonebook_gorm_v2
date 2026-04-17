package repository

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

// GetUsers
func (s *Service) GetUsers() ([]db.User, error) {
	var users []db.User
	err := s.DB.Preload("Phones").Find(&users).Error
	return users, err
}

// CreateUser
func (s *Service) CreateUser(user *db.User) error {
	return s.DB.Create(user).Error
}

// UpdateUser
func (s *Service) UpdateUser(user *db.User) error {
	return s.DB.Save(user).Error
}

// DeleteUser
func (s *Service) DeleteUser(userID uint) error {
	return s.DB.Delete(&db.User{}, userID).Error
}

// CreatePhone
func (s *Service) CreatePhone(phone *db.Phone) error {
	return s.DB.Create(phone).Error
}

// GetPhonesByUser
func (s *Service) GetPhonesByUser(userID uint) ([]db.Phone, error) {
	var phones []db.Phone
	err := s.DB.Where("user_id = ?", userID).Find(&phones).Error
	return phones, err
}

// DeletePhone
func (s *Service) DeletePhone(phoneID uint) error {
	return s.DB.Delete(&db.Phone{}, phoneID).Error
}

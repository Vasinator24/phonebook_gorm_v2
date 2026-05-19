package repository

import (
	"fmt"
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

	// clean IDs
	user.ID = 0

	for i := range user.Phones {
		user.Phones[i].ID = 0
		user.Phones[i].UserID = 0
	}

	// create user
	if err := s.DB.Create(user).Error; err != nil {
		return err
	}

	// if no phones → OK (no error)
	if len(user.Phones) == 0 {
		return nil
	}

	// assign FK
	for i := range user.Phones {
		user.Phones[i].UserID = user.ID
	}

	// try insert phones BUT don't break user creation
	if err := s.DB.Create(&user.Phones).Error; err != nil {
		// log error but DO NOT fail request
		fmt.Println("phones insert failed:", err)
	}

	return nil
}

// UpdateUser
func (s *Service) UpdateUser(user *db.User) error {
	return s.DB.Save(user).Error
}

// DeleteUser
func (s *Service) DeleteUser(userID uint) error {
	return s.DB.Delete(&db.User{}, userID).Error
}

func (s *Service) GetUserByEmail(email string) (*db.User, error) {
	var user db.User

	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreatePhone
func (s *Service) CreatePhone(phone *db.Phone) error {
	return s.DB.Create(phone).Error
}

// editPhone
func (s *Service) UpdatePhone(phone *db.Phone) error {
	return s.DB.Save(phone).Error
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

func (s *Service) GetPhoneByID(id uint) (*db.Phone, error) {
	var phone db.Phone

	err := s.DB.First(&phone, id).Error
	if err != nil {
		return nil, err
	}

	return &phone, nil
}

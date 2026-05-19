package services

import (
	"phonebook_gorm/db"
)

// CreatePhone създава нов телефонен номер
func (s *UserService) CreatePhone(phone *db.Phone) error {
	return s.repo.CreatePhone(phone)
}

// GetPhonesByUser връща всички телефони за даден user
func (s *UserService) GetPhonesByUser(userID uint) ([]db.Phone, error) {
	return s.repo.GetPhonesByUser(userID)
}

// DeletePhone изтрива телефон по ID
func (s *UserService) DeletePhone(phoneID uint) error {
	return s.repo.DeletePhone(phoneID)
}

func (s *UserService) GetPhoneByID(id uint) (*db.Phone, error) {
	return s.repo.GetPhoneByID(id)
}
func (s *UserService) UpdatePhone(phone *db.Phone) error {
	return s.repo.UpdatePhone(phone)
}

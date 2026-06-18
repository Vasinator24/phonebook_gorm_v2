package services

import "phonebook_gorm/db"

func (s *UserService) CreatePhone(phone *db.Phone) error {
	return s.repo.CreatePhone(phone)
}

func (s *UserService) GetPhonesByUser(userID uint) ([]db.Phone, error) {
	return s.repo.GetPhonesByUser(userID)
}

func (s *UserService) GetPhones() ([]db.PhoneWithUser, error) {
	return s.repo.GetPhones()
}

func (s *UserService) PhoneExists(number string, excludeID uint) (bool, error) {
	return s.repo.PhoneExists(number, excludeID)
}

func (s *UserService) DeletePhone(phoneID uint) error {
	return s.repo.DeletePhone(phoneID)
}

func (s *UserService) GetPhoneByID(id uint) (*db.Phone, error) {
	return s.repo.GetPhoneByID(id)
}

func (s *UserService) UpdatePhone(phone *db.Phone) error {
	return s.repo.UpdatePhone(phone)
}

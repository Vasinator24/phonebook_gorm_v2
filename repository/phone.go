package repository

import "phonebook_gorm/db"

func (s *Service) CreatePhone(phone *db.Phone) error {
	return s.DB.Create(phone).Error
}

func (s *Service) UpdatePhone(phone *db.Phone) error {
	return s.DB.Model(&db.Phone{}).
		Where("id = ?", phone.ID).
		Updates(map[string]interface{}{
			"number":  phone.Number,
			"user_id": phone.UserID,
		}).Error
}

func (s *Service) GetPhonesByUser(userID uint) ([]db.Phone, error) {
	var phones []db.Phone
	err := s.DB.Where("user_id = ?", userID).Find(&phones).Error
	return phones, err
}

func (s *Service) GetPhones() ([]db.PhoneWithUser, error) {
	var phones []db.PhoneWithUser
	err := s.DB.
		Table("phones").
		Select("phones.id, phones.user_id, phones.number, users.names AS user_name").
		Joins("JOIN users ON users.id = phones.user_id").
		Scan(&phones).Error

	return phones, err
}

func (s *Service) PhoneExists(number string, excludeID uint) (bool, error) {
	var count int64
	query := s.DB.Model(&db.Phone{}).Where("number = ?", number)

	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

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

package repository

import (
	"fmt"
	"time"

	"phonebook_gorm/db"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(dbConn *gorm.DB) *Service {
	return &Service{DB: dbConn}
}

func (s *Service) GetUsers() ([]db.User, error) {
	var users []db.User
	err := s.DB.Preload("Phones").Find(&users).Error
	return users, err
}

func (s *Service) CreateUser(user *db.User) error {
	user.ID = 0

	for i := range user.Phones {
		user.Phones[i].ID = 0
		user.Phones[i].UserID = 0
	}

	if err := s.DB.Create(user).Error; err != nil {
		return err
	}

	if len(user.Phones) == 0 {
		return nil
	}

	for i := range user.Phones {
		user.Phones[i].UserID = user.ID
	}

	if err := s.DB.Create(&user.Phones).Error; err != nil {
		fmt.Println("phones insert failed:", err)
	}

	return nil
}

func (s *Service) UpdateUser(id string, names string, email string) error {
	return s.DB.Model(&db.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"names": names,
			"email": email,
		}).Error
}

func (s *Service) DeleteUser(userID uint) error {
	return s.DB.Delete(&db.User{}, userID).Error
}

func (s *Service) GetUserByEmail(email string) (*db.User, error) {
	var user db.User

	err := s.DB.Where("email = ? OR username = ?", email, email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) GetUserByID(id uint) (*db.User, error) {
	var user db.User

	err := s.DB.Preload("Phones").First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) BlacklistToken(tokenHash string, expiresAt time.Time) error {
	return s.DB.Create(&db.BlacklistedToken{
		TokenHash: tokenHash,
		ExpiresAt: expiresAt.Unix(),
	}).Error
}

func (s *Service) DeleteExpiredBlacklistedTokens() error {
	return s.DB.
		Where("expires_at <= ?", time.Now().Unix()).
		Delete(&db.BlacklistedToken{}).Error
}

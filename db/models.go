package db

type User struct {
	ID           uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string  `json:"username" gorm:"unique"`
	Email        string  `json:"email"`
	Names        string  `json:"names"`
	PasswordHash string  `json:"-" gorm:"column:password_hash"`
	Phones       []Phone `json:"phones" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Phone struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"userID"`
	Number string `json:"number"`
}

type PhoneWithUser struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"userID"`
	Number   string `json:"number"`
	UserName string `json:"userName"`
}

type BlacklistedToken struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	TokenHash string `gorm:"uniqueIndex;not null"`
	ExpiresAt int64  `gorm:"not null"`
}

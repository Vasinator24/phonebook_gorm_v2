package db

type User struct {
	ID       uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string  `json:"username" gorm:"unique"`
	Email    string  `json:"email"`
	Names    string  `json:"names"`
	Phones   []Phone `json:"phones" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Phone struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"userID"`
	Number string `json:"number"`
}

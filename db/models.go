package db

type User struct {
	ID       uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string  `json:"username" gorm:"unique"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Names    string  `json:"names"`
	Role     string  `json:"role" gorm:"default:user"`
	Phones   []Phone `json:"phones" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Phone struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"userID"`
	Number string `json:"number"`
}

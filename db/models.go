package db

type User struct {
	ID       uint    `json:"id"`
	Username string  `json:"username" gorm:"unique"`
	Email    string  `json:"email"`
	Password string  `json:"-"`
	Names    string  `json:"names"`
	Phones   []Phone `json:"phones"`
}

type Phone struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Number string `json:"number"`
}

package db

import (
	"log"

	"gorm.io/gorm"
)

func Seed(dbConn *gorm.DB) {

	log.Println("Seeding database...")

	var count int64
	dbConn.Model(&User{}).Count(&count)

	// ако има данни — не правим нищо
	if count > 0 {
		log.Println("Seed skipped (data already exists)")
		return
	}

	users := []User{
		{
			Username: "ivan",
			Email:    "ivan@mail.com",
			Password: "1234",
			Names:    "Ivan Ivanov",
			Phones: []Phone{
				{Number: "0888123456"},
				{Number: "0899123456"},
			},
		},
		{
			Username: "george",
			Email:    "geo@mail.com",
			Password: "1234",
			Names:    "George Petrov",
			Phones: []Phone{
				{Number: "0877123456"},
			},
		},
	}

	for _, u := range users {
		dbConn.Create(&u)
	}

	log.Println("Seeding done")
}

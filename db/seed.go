package db

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(dbConn *gorm.DB) {

	log.Println("Seeding database...")

	// RESET TABLES + AUTO INCREMENT
	dbConn.Exec("SET FOREIGN_KEY_CHECKS = 0")

	dbConn.Exec("TRUNCATE TABLE phones")
	dbConn.Exec("TRUNCATE TABLE users")

	dbConn.Exec("SET FOREIGN_KEY_CHECKS = 1")

	log.Println("Database truncated")

	// HASH PASSWORD
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("1234"),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Fatal("Could not hash password")
	}

	users := []User{
		{
			Username: "admin",
			Email:    "admin@mail.com",
			Password: string(hashedPassword),
			Names:    "Admin User",
			Role:     "admin",
		},
		{
			Username: "ivan",
			Email:    "ivan@mail.com",
			Password: string(hashedPassword),
			Names:    "Ivan Ivanov",
			Phones: []Phone{
				{Number: "0888123456"},
				{Number: "0899123456"},
			},
		},
		{
			Username: "george",
			Email:    "geo@mail.com",
			Password: string(hashedPassword),
			Names:    "George Petrov",
			Phones: []Phone{
				{Number: "0877123456"},
			},
		},
		{
			Username: "maria",
			Email:    "maria@mail.com",
			Password: string(hashedPassword),
			Names:    "Maria Ivanova",
		},
		{
			Username: "nikolay",
			Email:    "nik@mail.com",
			Password: string(hashedPassword),
			Names:    "Nikolay Dimitrov",
			Phones: []Phone{
				{Number: "0888000001"},
			},
		},
		{
			Username: "stefan",
			Email:    "stefan@mail.com",
			Password: string(hashedPassword),
			Names:    "Stefan Kolev",
			Phones: []Phone{
				{Number: "0888000002"},
				{Number: "0888000003"},
				{Number: "0888000004"},
			},
		},
		{
			Username: "petya",
			Email:    "petya@mail.com",
			Password: string(hashedPassword),
			Names:    "Petya Petrova",
		},
		{
			Username: "alex",
			Email:    "alex@mail.com",
			Password: string(hashedPassword),
			Names:    "Alex Georgiev",
			Phones: []Phone{
				{Number: "0888000005"},
			},
		},
		{
			Username: "dani",
			Email:    "dani@mail.com",
			Password: string(hashedPassword),
			Names:    "Daniel Ivanov",
			Phones: []Phone{
				{Number: "0888000006"},
				{Number: "0888000007"},
			},
		},
		{
			Username: "teodora",
			Email:    "teo@mail.com",
			Password: string(hashedPassword),
			Names:    "Teodora Hristova",
		},
		{
			Username: "martin",
			Email:    "martin@mail.com",
			Password: string(hashedPassword),
			Names:    "Martin Petrov",
			Phones: []Phone{
				{Number: "0888000008"},
			},
		},
		{
			Username: "boris",
			Email:    "boris@mail.com",
			Password: string(hashedPassword),
			Names:    "Boris Nikolov",
		},
		{
			Username: "viktoria",
			Email:    "viki@mail.com",
			Password: string(hashedPassword),
			Names:    "Viktoria Dimitrova",
			Phones: []Phone{
				{Number: "0888000009"},
				{Number: "0888000010"},
			},
		},
		{
			Username: "emil",
			Email:    "emil@mail.com",
			Password: string(hashedPassword),
			Names:    "Emil Stoyanov",
		},
		{
			Username: "kalina",
			Email:    "kalina@mail.com",
			Password: string(hashedPassword),
			Names:    "Kalina Georgieva",
			Phones: []Phone{
				{Number: "0888000011"},
			},
		},
		{
			Username: "radi",
			Email:    "radi@mail.com",
			Password: string(hashedPassword),
			Names:    "Radoslav Kolev",
			Phones: []Phone{
				{Number: "0888000012"},
				{Number: "0888000013"},
				{Number: "0888000014"},
			},
		},
		{
			Username: "ani",
			Email:    "ani@mail.com",
			Password: string(hashedPassword),
			Names:    "Ani Petrova",
		},
		{
			Username: "kristian",
			Email:    "kris@mail.com",
			Password: string(hashedPassword),
			Names:    "Kristian Ivanov",
			Phones: []Phone{
				{Number: "0888000015"},
			},
		},
		{
			Username: "simona",
			Email:    "simona@mail.com",
			Password: string(hashedPassword),
			Names:    "Simona Angelova",
			Phones: []Phone{
				{Number: "0888000016"},
				{Number: "0888000017"},
			},
		},
		{
			Username: "ivo",
			Email:    "ivo@mail.com",
			Password: string(hashedPassword),
			Names:    "Ivo Dimitrov",
		},
		{
			Username: "plamen",
			Email:    "plamen@mail.com",
			Password: string(hashedPassword),
			Names:    "Plamen Georgiev",
			Phones: []Phone{
				{Number: "0888000018"},
			},
		},
	}

	for _, u := range users {
		dbConn.Create(&u)
	}

	log.Println("Seeding done")
}

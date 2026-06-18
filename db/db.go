package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("DB_DSN")

	log.Println("Connecting to DB...")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{}, &Phone{}, &BlacklistedToken{})

	log.Println("DB connected successfully")

	return db, nil
}

func Migrate(dbConn *gorm.DB) {
	log.Println("Running migrations...")

	err := dbConn.AutoMigrate(&User{}, &Phone{}, &BlacklistedToken{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migrations done")
}

func Reset(dbConn *gorm.DB) {
	log.Println("Resetting database...")

	// изтриваме в правилен ред
	dbConn.Exec("DELETE FROM phones")
	dbConn.Exec("DELETE FROM users")

	log.Println("Database cleared")
}

package main

import (
	"log"
	"os"

	controller "phonebook_gorm/controler"
	"phonebook_gorm/db"
	"phonebook_gorm/logger"
	"phonebook_gorm/repository"
	server "phonebook_gorm/server"
	"phonebook_gorm/services"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func runMigrate() {
	database, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	db.Migrate(database)
}

func runSeed() {
	database, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	db.Seed(database)
}

func runReset() {
	database, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	db.Reset(database)
}

func runServer() {
	app := fx.New(
		fx.Provide(
			db.NewDB,
			repository.NewService,
			services.NewUserService,
			controller.NewUserController,
			controller.NewPhoneController,
			logger.NewLogger,
		),
		fx.Invoke(
			server.StartServer,
		),
	)

	app.Run()
}
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	action := "server"

	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	switch action {

	case "migrate":
		runMigrate()

	case "seed":
		runSeed()

	case "reset":
		runReset()

	case "server":
		runServer()

	default:
		log.Println("Use: server | migrate | seed | reset")
	}

}

package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	controller "phonebook_gorm/controler"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

func getServerAddr() string {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	if strings.HasPrefix(port, ":") {
		return port
	}

	return ":" + port
}

func StartServer(
	lc fx.Lifecycle,
	userCtrl *controller.UserController,
	phoneCtrl *controller.PhoneController,
	dbConn *gorm.DB,
) {

	mux := http.NewServeMux()
	RegisterRoutes(mux, userCtrl, phoneCtrl, dbConn)
	addr := getServerAddr()

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal(err)
				}
			}()
			log.Println("Server running on " + addr)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server...")
			return server.Shutdown(ctx)
		},
	})
}

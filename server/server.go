package server

import (
	"context"
	"log"
	"net/http"
	controller "phonebook_gorm/controler"

	"go.uber.org/fx"
)

func StartServer(lc fx.Lifecycle, userCtrl *controller.UserController) {

	mux := http.NewServeMux()
	RegisterRoutes(mux, userCtrl)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal(err)
				}
			}()
			log.Println("Server running on :8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server...")
			return server.Shutdown(ctx)
		},
	})
}

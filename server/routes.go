package server

import (
	"net/http"

	"phonebook_gorm/auth"
	controller "phonebook_gorm/controler"

	"gorm.io/gorm"
)

func RegisterRoutes(
	mux *http.ServeMux,
	userCtrl *controller.UserController,
	phoneCtrl *controller.PhoneController,
	dbConn *gorm.DB,
) {
	protected := auth.AuthMiddleware(dbConn)

	// AUTH
	mux.Handle("/login",
		Cors(http.HandlerFunc(userCtrl.Login)),
	)

	mux.Handle("/logout",
		Cors(protected(http.HandlerFunc(userCtrl.Logout))),
	)

	mux.Handle("/me",
		Cors(protected(http.HandlerFunc(userCtrl.Me))),
	)

	// USERS
	mux.Handle("/users",
		Cors(protected(http.HandlerFunc(userCtrl.GetUsers))),
	)

	mux.Handle("/users/create",
		Cors(protected(http.HandlerFunc(userCtrl.CreateUser))),
	)

	mux.Handle("/users/delete",
		Cors(protected(http.HandlerFunc(userCtrl.DeleteUser))),
	)

	mux.Handle("/users/update",
		Cors(protected(http.HandlerFunc(userCtrl.UpdateUser))),
	)

	// PHONES
	mux.Handle("/phones",
		Cors(protected(http.HandlerFunc(phoneCtrl.GetPhonesByUser))),
	)

	mux.Handle("/phones/create",
		Cors(protected(http.HandlerFunc(phoneCtrl.CreatePhone))),
	)

	mux.Handle("/phones/delete",
		Cors(protected(http.HandlerFunc(phoneCtrl.DeletePhone))),
	)

	mux.Handle("/phones/update",
		Cors(protected(http.HandlerFunc(phoneCtrl.UpdatePhone))),
	)

	// HEALTH
	mux.Handle("/health",
		Cors(protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))),
	)
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

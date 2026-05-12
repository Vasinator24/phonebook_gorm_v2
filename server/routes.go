package server

import (
	"net/http"

	"phonebook_gorm/auth"
	controller "phonebook_gorm/controler"
)

func RegisterRoutes(mux *http.ServeMux, userCtrl *controller.UserController) {

	phoneCtrl := controller.NewPhoneController(userCtrl.GetService())

	// LOGIN
	mux.HandleFunc("/login", Cors(userCtrl.Login))

	// USERS
	mux.HandleFunc("/users", Cors(auth.AuthMiddleware(userCtrl.GetUsers)))
	mux.HandleFunc("/users/create", Cors(userCtrl.CreateUser))
	mux.HandleFunc("/users/delete",
		Cors(auth.AuthMiddleware(auth.AdminOnly(userCtrl.DeleteUser))),
	)

	mux.HandleFunc("/users/update",
		Cors(auth.AuthMiddleware(auth.AdminOnly(userCtrl.UpdateUser))),
	)

	// PHONES
	mux.HandleFunc("/phones", Cors(auth.AuthMiddleware(phoneCtrl.GetPhonesByUser)))
	mux.HandleFunc("/phones/create", Cors(auth.AuthMiddleware(phoneCtrl.CreatePhone)))
	mux.HandleFunc("/phones/delete", Cors(auth.AuthMiddleware(phoneCtrl.DeletePhone)))

	// HEALTH
	mux.HandleFunc("/health", Cors(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
}

func Cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

package server

import (
	"net/http"

	controller "phonebook_gorm/controler"
)

func RegisterRoutes(
	mux *http.ServeMux,
	userCtrl *controller.UserController,
	phoneCtrl *controller.PhoneController,
) {
	// USERS
	mux.Handle("/users",
		Cors(http.HandlerFunc(userCtrl.GetUsers)),
	)

	mux.Handle("/users/create",
		Cors(http.HandlerFunc(userCtrl.CreateUser)),
	)

	mux.Handle("/users/delete",
		Cors(http.HandlerFunc(userCtrl.DeleteUser)),
	)

	mux.Handle("/users/update",
		Cors(http.HandlerFunc(userCtrl.UpdateUser)),
	)

	// PHONES
	mux.Handle("/phones",
		Cors(http.HandlerFunc(phoneCtrl.GetPhonesByUser)),
	)

	mux.Handle("/phones/create",
		Cors(http.HandlerFunc(phoneCtrl.CreatePhone)),
	)

	mux.Handle("/phones/delete",
		Cors(http.HandlerFunc(phoneCtrl.DeletePhone)),
	)

	mux.Handle("/phones/update",
		Cors(http.HandlerFunc(phoneCtrl.UpdatePhone)),
	)

	// HEALTH
	mux.Handle("/health",
		Cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})),
	)
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

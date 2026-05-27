package server

import (
	"net/http"

	"phonebook_gorm/auth"
	controller "phonebook_gorm/controler"
)

func RegisterRoutes(mux *http.ServeMux, userCtrl *controller.UserController) {

	phoneCtrl := controller.NewPhoneController(userCtrl.GetService())

	// LOGIN
	mux.Handle("/login",
		Cors(http.HandlerFunc(userCtrl.Login)),
	)

	// USERS
	mux.Handle("/users",
		Cors(auth.AuthMiddleware(http.HandlerFunc(userCtrl.GetUsers))),
	)

	mux.Handle("/users/create",
		Cors(http.HandlerFunc(userCtrl.CreateUser)),
	)

	mux.Handle("/users/admin-create",
		Cors(auth.AuthMiddleware(auth.AdminOnly(http.HandlerFunc(userCtrl.AdminCreateUser)))),
	)

	mux.Handle("/users/delete",
		Cors(auth.AuthMiddleware(auth.AdminOnly(http.HandlerFunc(userCtrl.DeleteUser)))),
	)

	mux.Handle("/users/update",
		Cors(auth.AuthMiddleware(http.HandlerFunc(userCtrl.UpdateUser))),
	)

	// PHONES
	mux.Handle("/phones",
		Cors(auth.AuthMiddleware(http.HandlerFunc(phoneCtrl.GetPhonesByUser))),
	)

	mux.Handle("/phones/create",
		Cors(auth.AuthMiddleware(http.HandlerFunc(phoneCtrl.CreatePhone))),
	)

	mux.Handle("/phones/delete",
		Cors(auth.AuthMiddleware(http.HandlerFunc(phoneCtrl.DeletePhone))),
	)

	mux.Handle("/phones/update",
		Cors(auth.AuthMiddleware(http.HandlerFunc(phoneCtrl.UpdatePhone))),
	)

	// HEALTH
	mux.Handle("/health",
		Cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

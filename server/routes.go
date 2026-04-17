package server

import (
	"net/http"

	"phonebook_gorm/auth"
	controller "phonebook_gorm/controler"
)

func RegisterRoutes(mux *http.ServeMux, userCtrl *controller.UserController) {

	mux.HandleFunc("/users", userCtrl.GetUsers)
	//mux.HandleFunc("/users/create", userCtrl.CreateUser)
	//mux.HandleFunc("/users/update", userCtrl.UpdateUser)
	//mux.HandleFunc("/users/delete", userCtrl.DeleteUser)

	mux.HandleFunc("/login", userCtrl.Login)
	// protected route
	//mux.HandleFunc("/users", auth.AuthMiddleware(userCtrl.GetUsers))
	mux.HandleFunc("/users/create", auth.AuthMiddleware(userCtrl.CreateUser))
	mux.HandleFunc("/users/update", auth.AuthMiddleware(userCtrl.UpdateUser))
	mux.HandleFunc("/users/delete", auth.AuthMiddleware(userCtrl.DeleteUser))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}

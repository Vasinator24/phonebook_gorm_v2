package server

import (
	"net/http"

	controller "phonebook_gorm/controler"
)

func RegisterRoutes(mux *http.ServeMux, userCtrl *controller.UserController) {

	mux.HandleFunc("/api/users", userCtrl.GetUsers)
	mux.HandleFunc("/api/users/create", userCtrl.CreateUser)

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}

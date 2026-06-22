package server

import (
	"net/http"

	"phonebook_gorm/auth"
	controller "phonebook_gorm/controler"
)

const APIV1 = "/api/v1"

func RegisterRoutes(
	mux *http.ServeMux,
	userCtrl *controller.UserController,
	phoneCtrl *controller.PhoneController,
	authServe *auth.AuthServe,
) {
	apiV1 := NewRouteGroup(mux, APIV1, NewCorsMiddleware())
	protectedAPIv1 := apiV1.With(authServe.Middleware)

	registerAuthRoutes(apiV1, protectedAPIv1, userCtrl)
	registerUserRoutes(protectedAPIv1, userCtrl)
	registerPhoneRoutes(protectedAPIv1, phoneCtrl)

	protectedAPIv1.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
}

func registerAuthRoutes(
	public *RouteGroup,
	protected *RouteGroup,
	userCtrl *controller.UserController,
) {
	public.Handle("/login", http.HandlerFunc(userCtrl.Login))
	protected.Handle("/logout", http.HandlerFunc(userCtrl.Logout))
	protected.Handle("/me", http.HandlerFunc(userCtrl.Me))
}

func registerUserRoutes(group *RouteGroup, userCtrl *controller.UserController) {
	group.Handle("/users", http.HandlerFunc(userCtrl.GetUsers))
	group.Handle("/users/create", http.HandlerFunc(userCtrl.CreateUser))
	group.Handle("/users/delete", http.HandlerFunc(userCtrl.DeleteUser))
	group.Handle("/users/update", http.HandlerFunc(userCtrl.UpdateUser))
}

func registerPhoneRoutes(group *RouteGroup, phoneCtrl *controller.PhoneController) {
	group.Handle("/phones", http.HandlerFunc(phoneCtrl.GetPhonesByUser))
	group.Handle("/phones/create", http.HandlerFunc(phoneCtrl.CreatePhone))
	group.Handle("/phones/delete", http.HandlerFunc(phoneCtrl.DeletePhone))
	group.Handle("/phones/update", http.HandlerFunc(phoneCtrl.UpdatePhone))
}

package controller

import (
	"encoding/json"
	"net/http"

	"phonebook_gorm/db"
	"phonebook_gorm/logger"
	"phonebook_gorm/services"
)

type UserController struct {
	service *services.UserService
	log     *logger.Logger
}

func NewUserController(s *services.UserService, l *logger.Logger) *UserController {
	return &UserController{
		service: s,
		log:     l,
	}
}

// GET /api/users
func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	c.log.Info.Info("GetUsers called")

	users, err := c.service.GetUsers()
	if err != nil {
		c.log.Error.Error("DB error while fetching users")
		http.Error(w, err.Error(), 500)
		return
	}

	c.log.Warn.Info("Users successfully fetched")

	json.NewEncoder(w).Encode(users)

}

// POST /api/users/create
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {

	c.log.Debug.Info("CreateUser request received")

	var user db.User
	json.NewDecoder(r.Body).Decode(&user)

	c.log.Info.Info("Creating user started")

	err := c.service.CreateUser(&user)
	if err != nil {
		c.log.Error.Error("failed to create user")
		http.Error(w, err.Error(), 500)
		return
	}

	c.log.Info.Info("User created successfully")

	json.NewEncoder(w).Encode(user)
}

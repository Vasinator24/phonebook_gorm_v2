package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"phonebook_gorm/auth"
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

// GET /users
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

// POST /users/create
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

// PUT /users/update
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	c.log.Info.Info("UpdateUser called")

	var user db.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.log.Error.Error("invalid request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.UpdateUser(&user)
	if err != nil {
		c.log.Error.Error("failed to update user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.log.Info.Info("user updated successfully")

	json.NewEncoder(w).Encode(user)
}

// DELETE /users/delete
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	c.log.Info.Info("DeleteUser called")

	// взимаме ID от query: /users?id=1
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.log.Error.Error("invalid user id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = c.service.DeleteUser(uint(id))
	if err != nil {
		c.log.Error.Error("failed to delete user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.log.Info.Info("user deleted successfully")

	w.WriteHeader(http.StatusNoContent)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {

	var input db.User

	json.NewDecoder(r.Body).Decode(&input)

	users, err := c.service.GetUsers()
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	for _, u := range users {
		if u.Email == input.Email {
			token, _ := auth.GenerateToken(u.ID, u.Email)

			json.NewEncoder(w).Encode(map[string]string{
				"token": token,
			})
			return
		}
	}

	http.Error(w, "invalid credentials", http.StatusUnauthorized)
}

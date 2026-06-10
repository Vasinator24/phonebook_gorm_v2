package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

}

// POST /users/create
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	c.log.Debug.Info("CreateUser request received")

	var user db.User
	json.NewDecoder(r.Body).Decode(&user)

	user.Username = strings.TrimSpace(user.Username)
	user.Names = strings.TrimSpace(user.Names)
	user.Email = strings.TrimSpace(user.Email)

	if len(user.Username) < 3 {
		http.Error(w, "username must be at least 3 characters", http.StatusBadRequest)
		return
	}

	if len(user.Names) < 3 {
		http.Error(w, "name must be at least 3 characters", http.StatusBadRequest)
		return
	}

	if !strings.Contains(user.Email, "@") {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	c.log.Info.Info("Creating user started")

	err := c.service.CreateUser(&user)
	if err != nil {
		c.log.Error.Error("failed to create user")
		http.Error(w, err.Error(), 400)
		return
	}

	c.log.Info.Info("User created successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// PUT /users/update
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	c.log.Info.Info("UpdateUser called")

	// get ID from query
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// decode body
	var user db.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.log.Error.Error("invalid request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = uint(id)
	user.Names = strings.TrimSpace(user.Names)
	user.Email = strings.TrimSpace(user.Email)

	if len(user.Names) < 3 {
		http.Error(w, "name must be at least 3 characters", http.StatusBadRequest)
		return
	}

	if !strings.Contains(user.Email, "@") {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	// update
	err = c.service.UpdateUser(idStr, user.Names, user.Email)
	if err != nil {
		c.log.Error.Error("failed to update user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.log.Info.Info("user updated successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DELETE /users/delete
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	c.log.Info.Info("DeleteUser called")

	// get ID
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.log.Error.Error("invalid user id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// delete
	err = c.service.DeleteUser(uint(id))
	if err != nil {
		c.log.Error.Error("failed to delete user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.log.Info.Info("user deleted successfully")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user deleted"))
}

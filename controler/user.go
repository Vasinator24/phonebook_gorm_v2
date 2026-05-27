package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"phonebook_gorm/auth"
	"phonebook_gorm/db"
	"phonebook_gorm/logger"
	"phonebook_gorm/services"

	"golang.org/x/crypto/bcrypt"
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

func (uc *UserController) GetService() *services.UserService {
	return uc.service
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// четем body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// намираме user
	user, err := c.service.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	// проверка на password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	// генерираме token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	// ВРЪЩАМЕ TOKEN
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
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

func (c *UserController) createUser(w http.ResponseWriter, r *http.Request, allowRole bool) {
	c.log.Debug.Info("CreateUser request received")

	var user db.User
	json.NewDecoder(r.Body).Decode(&user)

	if !allowRole || user.Role == "" {
		user.Role = "user"
	}

	c.log.Info.Info("Creating user started")

	err := c.service.CreateUser(&user)
	if err != nil {
		c.log.Error.Error("failed to create user")
		http.Error(w, err.Error(), 400)
		return
	}

	c.log.Info.Info("User created successfully")

	json.NewEncoder(w).Encode(user)
}

// POST /users/create
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	c.createUser(w, r, false)
}

// POST /users/admin-create
func (c *UserController) AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	c.createUser(w, r, true)
}

// PUT /users/update
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	c.log.Info.Info("UpdateUser called")

	// AUTH
	claims := auth.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// get ID from query
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// AUTHORIZATION RULE
	if claims.Role != "admin" && claims.UserID != uint(id) {
		http.Error(w, "forbidden", http.StatusForbidden)
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

	// update
	err = c.service.UpdateUser(idStr, user.Names, user.Email)
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

	// AUTHORIZATION
	claims := auth.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// get ID
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.log.Error.Error("invalid user id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// AUTH RULE
	if claims.Role != "admin" && claims.UserID != uint(id) {
		http.Error(w, "forbidden", http.StatusForbidden)
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

	w.WriteHeader(http.StatusNoContent)
}

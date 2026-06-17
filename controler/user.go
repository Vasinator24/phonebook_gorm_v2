package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Names    string `json:"names"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user := db.User{
		Username: strings.TrimSpace(req.Username),
		Names:    strings.TrimSpace(req.Names),
		Email:    strings.TrimSpace(req.Email),
	}

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

	if len(req.Password) < 4 {
		http.Error(w, "password must be at least 4 characters", http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.log.Error.Error("failed to hash password")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.PasswordHash = string(passwordHash)

	c.log.Info.Info("Creating user started")

	err = c.service.CreateUser(&user)
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

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	c.log.Info.Info("Login called")

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUserByEmail(strings.TrimSpace(req.Email))
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.log.Error.Error("failed to generate token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	auth.SetTokenCookie(w, token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	if err := c.service.DeleteExpiredBlacklistedTokens(); err != nil {
		c.log.Error.Error("failed to delete expired blacklisted tokens")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie(auth.CookieName)
	if err == nil {
		claims, parseErr := auth.ParseToken(cookie.Value)
		if parseErr == nil && claims.ExpiresAt != nil {
			tokenHash := auth.TokenHash(cookie.Value)
			if err := c.service.BlacklistToken(tokenHash, claims.ExpiresAt.Time); err != nil {
				c.log.Error.Error("failed to blacklist token")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	auth.ClearTokenCookie(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("logged out"))
}

func (c *UserController) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := c.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

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

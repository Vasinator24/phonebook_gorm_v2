package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"phonebook_gorm/auth"
	"phonebook_gorm/db"
	"phonebook_gorm/services"
)

type PhoneController struct {
	service *services.UserService
}

func NewPhoneController(s *services.UserService) *PhoneController {
	return &PhoneController{service: s}
}

// Create Phone
func (pc *PhoneController) CreatePhone(w http.ResponseWriter, r *http.Request) {

	// AUTH
	claims := auth.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var phone db.Phone

	if err := json.NewDecoder(r.Body).Decode(&phone); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// AUTHORIZATION RULE
	if claims.Role != "admin" && claims.UserID != phone.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err := pc.service.CreatePhone(&phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(phone)
}

// Get Phones by User
func (pc *PhoneController) GetPhonesByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	phones, err := pc.service.GetPhonesByUser(uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(phones)
}

// Delete Phone
func (pc *PhoneController) DeletePhone(w http.ResponseWriter, r *http.Request) {

	// AUTH
	claims := auth.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	// взимаме телефона от DB
	phone, err := pc.service.GetPhoneByID(uint(id))
	if err != nil {
		http.Error(w, "phone not found", http.StatusNotFound)
		return
	}

	// AUTHORIZATION RULE
	if claims.Role != "admin" && claims.UserID != phone.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// delete
	err = pc.service.DeletePhone(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Phone deleted"))
}

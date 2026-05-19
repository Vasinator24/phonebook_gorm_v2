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

// CREATE PHONE

func (pc *PhoneController) CreatePhone(w http.ResponseWriter, r *http.Request) {

	// AUTH
	claims := auth.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// REQUEST BODY
	var req struct {
		Number string `json:"number"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// CREATE PHONE FROM TOKEN USER
	phone := db.Phone{
		UserID: claims.UserID,
		Number: req.Number,
	}

	err = pc.service.CreatePhone(&phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(phone)
}

// UPDATE PHONE

func (pc *PhoneController) UpdatePhone(w http.ResponseWriter, r *http.Request) {

	claims := auth.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// GET PHONE FROM DB
	phone, err := pc.service.GetPhoneByID(uint(id))
	if err != nil {
		http.Error(w, "phone not found", http.StatusNotFound)
		return
	}

	// AUTHORIZATION
	if claims.Role != "admin" && claims.UserID != phone.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// REQUEST BODY
	var req struct {
		Number string `json:"number"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// UPDATE
	phone.Number = req.Number

	err = pc.service.UpdatePhone(phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(phone)
}

// GET PHONES BY USER

func (pc *PhoneController) GetPhonesByUser(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.URL.Query().Get("user_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	phones, err := pc.service.GetPhonesByUser(uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(phones)
}

// DELETE PHONE

func (pc *PhoneController) DeletePhone(w http.ResponseWriter, r *http.Request) {

	claims := auth.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// GET PHONE
	phone, err := pc.service.GetPhoneByID(uint(id))
	if err != nil {
		http.Error(w, "phone not found", http.StatusNotFound)
		return
	}

	// AUTHORIZATION
	if claims.Role != "admin" && claims.UserID != phone.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// DELETE
	err = pc.service.DeletePhone(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("phone deleted"))
}

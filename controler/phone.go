package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	// REQUEST BODY
	var req struct {
		UserID uint   `json:"user_id"`
		Number string `json:"number"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// CREATE PHONE FOR TARGET USER
	phone := db.Phone{
		UserID: req.UserID,
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

	// REQUEST BODY
	var req struct {
		UserID uint   `json:"user_id"`
		Number string `json:"number"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// UPDATE
	phone.UserID = req.UserID
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

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// GET PHONE
	_, err = pc.service.GetPhoneByID(uint(id))
	if err != nil {
		http.Error(w, "phone not found", http.StatusNotFound)
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

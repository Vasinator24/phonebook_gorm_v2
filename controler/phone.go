package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"phonebook_gorm/db"
	"phonebook_gorm/logger"
	"phonebook_gorm/services"
)

type PhoneController struct {
	service *services.UserService
	log     *logger.Logger
}

func NewPhoneController(s *services.UserService, l *logger.Logger) *PhoneController {
	return &PhoneController{
		service: s,
		log:     l,
	}
}

// CREATE PHONE

func (pc *PhoneController) CreatePhone(w http.ResponseWriter, r *http.Request) {
	pc.log.Info.Info("CreatePhone called")

	// REQUEST BODY
	var req struct {
		UserID uint   `json:"user_id"`
		Number string `json:"number"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		pc.log.Error.Error("invalid create phone request body")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 {
		pc.log.Error.Error("create phone failed: user_id is required")
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
		pc.log.Error.Error("failed to create phone")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pc.log.Info.Info("phone created successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(phone)
}

// UPDATE PHONE

func (pc *PhoneController) UpdatePhone(w http.ResponseWriter, r *http.Request) {
	pc.log.Info.Info("UpdatePhone called")

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		pc.log.Error.Error("invalid phone id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// GET PHONE FROM DB
	phone, err := pc.service.GetPhoneByID(uint(id))
	if err != nil {
		pc.log.Error.Error("phone not found")
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
		pc.log.Error.Error("invalid update phone request body")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 {
		pc.log.Error.Error("update phone failed: user_id is required")
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// UPDATE
	phone.UserID = req.UserID
	phone.Number = req.Number

	err = pc.service.UpdatePhone(phone)
	if err != nil {
		pc.log.Error.Error("failed to update phone")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pc.log.Info.Info("phone updated successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(phone)
}

// GET PHONES BY USER

func (pc *PhoneController) GetPhonesByUser(w http.ResponseWriter, r *http.Request) {
	pc.log.Info.Info("GetPhonesByUser called")

	userIDStr := r.URL.Query().Get("user_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		pc.log.Error.Error("invalid user_id while fetching phones")
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	phones, err := pc.service.GetPhonesByUser(uint(userID))
	if err != nil {
		pc.log.Error.Error("failed to fetch phones by user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pc.log.Info.Info("phones fetched successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(phones)
}

// DELETE PHONE

func (pc *PhoneController) DeletePhone(w http.ResponseWriter, r *http.Request) {
	pc.log.Info.Info("DeletePhone called")

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		pc.log.Error.Error("invalid phone id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// GET PHONE
	_, err = pc.service.GetPhoneByID(uint(id))
	if err != nil {
		pc.log.Error.Error("phone not found")
		http.Error(w, "phone not found", http.StatusNotFound)
		return
	}

	// DELETE
	err = pc.service.DeletePhone(uint(id))
	if err != nil {
		pc.log.Error.Error("failed to delete phone")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pc.log.Info.Info("phone deleted successfully")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("phone deleted"))
}

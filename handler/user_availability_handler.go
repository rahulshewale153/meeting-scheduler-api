package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/rahulshewale153/meeting-scheduler-api/service"
	"github.com/rahulshewale153/meeting-scheduler-api/utils"
)

type UserAvailabilityHandler struct {
	userAvailabilityService service.UserAvailabilityServiceI
}

func NewUserAvailabilityHandler(userAvailabilityService service.UserAvailabilityServiceI) *UserAvailabilityHandler {
	return &UserAvailabilityHandler{
		userAvailabilityService: userAvailabilityService,
	}
}

// InsertUserAvailability inserts a new user availability record
func (h *UserAvailabilityHandler) InsertUserAvailability(w http.ResponseWriter, r *http.Request) {
	var userAvailability model.UserAvailability
	err := json.NewDecoder(r.Body).Decode(&userAvailability)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	eventIDStr := vars["event_id"]
	if eventIDStr == "" {
		http.Error(w, "event_id is required", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid event_id: %v", err), http.StatusBadRequest)
		return
	}

	userIDStr := vars["user_id"]
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user_id: %v", err), http.StatusBadRequest)
		return
	}

	userAvailability.UserID = userID
	userAvailability.EventID = eventID
	if errs, ok := utils.IsValid(userAvailability); !ok {
		log.Printf("Validation failed: %v", errs)
		http.Error(w, fmt.Sprintf("Validation failed: %v", errs.Error), http.StatusBadRequest)
		return
	}

	err = h.userAvailabilityService.InsertUserAvailability(r.Context(), userAvailability)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User availability inserted successfully"})
}

// UpdateUserAvailability updates the availability of a user for a specific event
func (h *UserAvailabilityHandler) UpdateUserAvailability(w http.ResponseWriter, r *http.Request) {
	var userAvailability model.UserAvailability
	err := json.NewDecoder(r.Body).Decode(&userAvailability)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	eventIDStr := vars["event_id"]
	if eventIDStr == "" {
		http.Error(w, "event_id is required", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid event_id: %v", err), http.StatusBadRequest)
		return
	}

	userIDStr := vars["user_id"]
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user_id: %v", err), http.StatusBadRequest)
		return
	}

	userAvailability.UserID = userID
	userAvailability.EventID = eventID
	if errs, ok := utils.IsValid(userAvailability); !ok {
		log.Printf("Validation failed: %v", errs)
		http.Error(w, fmt.Sprintf("Validation failed: %v", errs.Error), http.StatusBadRequest)
		return
	}

	err = h.userAvailabilityService.UpdateUserAvailability(r.Context(), userAvailability)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User availability updated successfully"})
}

// GetUserAvailability retrieves the availability of a specific user for a specific event
func (h *UserAvailabilityHandler) GetUserAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventIDStr := vars["event_id"]
	if eventIDStr == "" {
		http.Error(w, "event_id is required", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid event_id: %v", err), http.StatusBadRequest)
		return
	}

	userIDStr := vars["user_id"]
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user_id: %v", err), http.StatusBadRequest)
		return
	}

	slots, err := h.userAvailabilityService.GetUserAvailability(r.Context(), eventID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(slots)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(slots)
}

// DeleteUserAvailability deletes a user's availability for a specific event
func (h *UserAvailabilityHandler) DeleteUserAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventIDStr := vars["event_id"]
	if eventIDStr == "" {
		http.Error(w, "event_id is required", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid event_id: %v", err), http.StatusBadRequest)
		return
	}

	userIDStr := vars["user_id"]
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user_id: %v", err), http.StatusBadRequest)
		return
	}

	err = h.userAvailabilityService.DeleteUserAvailability(r.Context(), userID, eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User availability deleted successfully"})
}

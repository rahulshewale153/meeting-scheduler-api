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

type EventHandler struct {
	eventService service.EventServiceI
}

func NewEventHandler(eventService service.EventServiceI) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

// InsertEvent inserts a new event
func (h *EventHandler) InsertEvent(w http.ResponseWriter, r *http.Request) {
	var req model.EventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if errs, ok := utils.IsValid(req); !ok {
		log.Printf("Validation failed: %v", errs)
		http.Error(w, fmt.Sprintf("Validation failed:%v", errs.Error), http.StatusBadRequest)
		return
	}

	eventID, err := h.eventService.InsertEvent(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"event_id": eventID})
}

// UpdateEvent updates an existing event
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var req model.EventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if errs, ok := utils.IsValid(req); !ok {
		log.Printf("Validation failed: %v", errs)
		http.Error(w, fmt.Sprintf("Validation failed:%v", errs.Error), http.StatusBadRequest)
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
		http.Error(w, "Invalid event_id", http.StatusBadRequest)
		return
	}

	req.ID = eventID // Set the ID in the request to update the specific event

	err = h.eventService.UpdateEvent(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteEvent deletes an existing event
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventIDStr := vars["event_id"]
	if eventIDStr == "" {
		http.Error(w, "event_id is required", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid event_id", http.StatusBadRequest)
		return
	}

	err = h.eventService.DeleteEvent(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

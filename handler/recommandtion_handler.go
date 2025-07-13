package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rahulshewale153/meeting-scheduler-api/service"
)

type RecommendationHandler struct {
	RecommendationService service.RecommendationServiceI
}

func NewRecommendationHandler(recommendationService service.RecommendationServiceI) *RecommendationHandler {
	return &RecommendationHandler{
		RecommendationService: recommendationService,
	}
}

func (h *RecommendationHandler) GetRecommendedSlots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventIDStr := vars["event_id"]
	if eventIDStr == "" {
		http.Error(w, "event_id query parameter is required", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid event_id", http.StatusBadRequest)
		return
	}

	recommendedSlots, err := h.RecommendationService.GetRecommendedSlots(r.Context(), eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendedSlots)
}

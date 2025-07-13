package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	mockService "github.com/rahulshewale153/meeting-scheduler-api/mock/service"
	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/assert"
)

func TestRecommendationHandler_GetRecommendedSlots(t *testing.T) {
	MockRecommendationService := new(mockService.MockRecommendationService)
	recommendationHandler := NewRecommendationHandler(MockRecommendationService)

	t.Run("GetRecommendedSlots should return an error when event_id is not provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/1/recommendation", nil)
		w := httptest.NewRecorder()

		recommendationHandler.GetRecommendedSlots(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("GetRecommendedSlots should return an error when event_id is invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/1/recommendation", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "invalid"})
		w := httptest.NewRecorder()

		recommendationHandler.GetRecommendedSlots(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("GetRecommendedSlots should return an error when service returns an error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/1/recommendation", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		w := httptest.NewRecorder()

		MockRecommendationService.On("GetRecommendedSlots", req.Context(), int64(1)).Return([]model.SlotRecommendation{}, assert.AnError).Once()

		recommendationHandler.GetRecommendedSlots(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("GetRecommendedSlots should return recommended slots when service returns successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/1/recommendation", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		w := httptest.NewRecorder()

		expectedSlots := []model.SlotRecommendation{
			{
				Slot: model.EventSlot{
					StartTime: time.Date(2025, 07, 13, 10, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2025, 07, 13, 11, 0, 0, 0, time.UTC),
				},
				Available:   []int64{1, 2},
				Unavailable: []int64{3},
			},
		}
		MockRecommendationService.On("GetRecommendedSlots", req.Context(), int64(1)).Return(expectedSlots, nil).Once()

		recommendationHandler.GetRecommendedSlots(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var actualSlots []model.SlotRecommendation
		err := json.NewDecoder(w.Body).Decode(&actualSlots)
		assert.NoError(t, err)

		assert.Equal(t, expectedSlots, actualSlots)
	})

}

package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	mockService "github.com/rahulshewale153/meeting-scheduler-api/mock/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsertEvent(t *testing.T) {
	mockEventService := new(mockService.MockEventService)
	eventHandler := NewEventHandler(mockEventService)
	validRequest := `{"title": "Test Event", "organizer_id": 1, "duration_minutes": 60, "proposed_slots": [{"start_time": "2023-10-01T10:00:00Z", "end_time": "2023-10-01T11:00:00Z"}]}`

	t.Run("invalid JSON request, should return an error", func(t *testing.T) {
		invalidRequest := `\invalid_json`
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(invalidRequest))
		w := httptest.NewRecorder()

		eventHandler.InsertEvent(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error, should return internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(validRequest))
		w := httptest.NewRecorder()

		mockEventService.On("InsertEvent", req.Context(), mock.Anything).Return(int64(0), assert.AnError).Once()

		eventHandler.InsertEvent(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("valid request, should return event ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(validRequest))
		w := httptest.NewRecorder()

		mockEventService.On("InsertEvent", req.Context(), mock.Anything).Return(int64(1), nil).Once()

		eventHandler.InsertEvent(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"event_id":1`)
	})
}

func TestUpdateEvent(t *testing.T) {
	mockEventService := new(mockService.MockEventService)
	eventHandler := NewEventHandler(mockEventService)
	validRequest := `{"id": 1, "title": "Updated Event", "organizer_id": 1, "duration_minutes": 60, "proposed_slots": [{"start_time": "2023-10-01T10:00:00Z", "end_time": "2023-10-01T11:00:00Z"}]}`

	t.Run("invalid JSON request, should return an error", func(t *testing.T) {
		invalidRequest := `\invalid_json`
		req := httptest.NewRequest(http.MethodPut, "/events/1", strings.NewReader(invalidRequest))
		w := httptest.NewRecorder()

		eventHandler.UpdateEvent(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error, should return internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/events/1", strings.NewReader(validRequest))
		w := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		mockEventService.On("UpdateEvent", req.Context(), mock.Anything).Return(assert.AnError).Once()

		eventHandler.UpdateEvent(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("valid request, should return no content status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/events/1", strings.NewReader(validRequest))
		w := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		mockEventService.On("UpdateEvent", req.Context(), mock.Anything).Return(nil).Once()

		eventHandler.UpdateEvent(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

}

func TestDeleteEvent(t *testing.T) {
	mockEventService := new(mockService.MockEventService)
	eventHandler := NewEventHandler(mockEventService)

	t.Run("invalid event ID, should return bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "invalid"})
		w := httptest.NewRecorder()

		eventHandler.DeleteEvent(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error, should return internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		w := httptest.NewRecorder()

		mockEventService.On("DeleteEvent", req.Context(), int64(1)).Return(assert.AnError).Once()

		eventHandler.DeleteEvent(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("valid request, should return no content status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		w := httptest.NewRecorder()

		mockEventService.On("DeleteEvent", req.Context(), int64(1)).Return(nil).Once()

		eventHandler.DeleteEvent(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

}

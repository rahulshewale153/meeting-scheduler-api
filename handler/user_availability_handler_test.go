package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	mockService "github.com/rahulshewale153/meeting-scheduler-api/mock/service"
	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsertUserAvailability(t *testing.T) {
	mockUserAvailService := new(mockService.MockUserAvailabilityService)
	userAvailabilityHandler := NewUserAvailabilityHandler(mockUserAvailService)
	validRequest := `{"availability": [{"start_time": "2023-10-01T10:00:00Z", "end_time": "2023-10-01T11:00:00Z"}]}`

	t.Run("invalid JSON request, should return an error", func(t *testing.T) {
		invalidRequest := `\invalid_json`
		req := httptest.NewRequest(http.MethodPost, "/user_availability", strings.NewReader(invalidRequest))
		w := httptest.NewRecorder()

		userAvailabilityHandler.InsertUserAvailability(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing event_id in URL, should return bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events/1/availability/1", strings.NewReader(validRequest))
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		w := httptest.NewRecorder()
		userAvailabilityHandler.InsertUserAvailability(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "event_id is required")
	})

	t.Run("missing user_id in URL, should return bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events/1/availability/1", strings.NewReader(validRequest))

		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		w := httptest.NewRecorder()
		userAvailabilityHandler.InsertUserAvailability(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user_id is required")
	})

	t.Run("service error, should return internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events/1/availability/1", strings.NewReader(validRequest))
		req = mux.SetURLVars(req, map[string]string{"user_id": "1", "event_id": "1"})
		w := httptest.NewRecorder()
		mockUserAvailService.On("InsertUserAvailability", req.Context(), mock.Anything).Return(assert.AnError).Once()

		userAvailabilityHandler.InsertUserAvailability(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUserAvailService.AssertExpectations(t)
	})

	t.Run("successful insert, should return created status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events/1/availability/1", strings.NewReader(validRequest))
		req = mux.SetURLVars(req, map[string]string{"user_id": "1", "event_id": "1"})
		w := httptest.NewRecorder()
		mockUserAvailService.On("InsertUserAvailability", req.Context(), mock.Anything).Return(nil).Once()

		userAvailabilityHandler.InsertUserAvailability(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		mockUserAvailService.AssertExpectations(t)
		assert.Contains(t, w.Body.String(), `{"message":"User availability inserted successfully"}`)
	})

}

func TestUpdateUserAvailability(t *testing.T) {
	mockUserAvailService := new(mockService.MockUserAvailabilityService)
	userAvailabilityHandler := NewUserAvailabilityHandler(mockUserAvailService)
	validRequest := `{"availability": [{"start_time": "2023-10-01T10:00:00Z", "end_time": "2023-10-01T11:00:00Z"}]}`

	t.Run("invalid JSON request, should return an error", func(t *testing.T) {
		invalidRequest := `\invalid_json`
		req := httptest.NewRequest(http.MethodPut, "/events/1/availability/1", strings.NewReader(invalidRequest))
		w := httptest.NewRecorder()

		userAvailabilityHandler.UpdateUserAvailability(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing event_id in URL, should return bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/events/1/availability/1", strings.NewReader(validRequest))
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		w := httptest.NewRecorder()
		userAvailabilityHandler.UpdateUserAvailability(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "event_id is required")
	})

	t.Run("missing user_id in URL, should return bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/events/1/availability/1", strings.NewReader(validRequest))
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		w := httptest.NewRecorder()
		userAvailabilityHandler.UpdateUserAvailability(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user_id is required")
	})

	t.Run("service error, should return internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/events/1/availability/1", strings.NewReader(validRequest))
		req = mux.SetURLVars(req, map[string]string{"user_id": "1", "event_id": "1"})
		w := httptest.NewRecorder()
		mockUserAvailService.On("UpdateUserAvailability", req.Context(), mock.Anything).Return(assert.AnError).Once()

		userAvailabilityHandler.UpdateUserAvailability(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), assert.AnError.Error())
		mockUserAvailService.AssertExpectations(t)
	})

	t.Run("successful update, should return no content", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/events/1/availability/1", strings.NewReader(validRequest))
		req = mux.SetURLVars(req, map[string]string{"user_id": "1", "event_id": "1"})
		w := httptest.NewRecorder()
		mockUserAvailService.On("UpdateUserAvailability", req.Context(), mock.Anything).Return(nil).Once()

		userAvailabilityHandler.UpdateUserAvailability(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockUserAvailService.AssertExpectations(t)

	})

}

func TestGetUserAvailability(t *testing.T) {
	mockUserAvailService := new(mockService.MockUserAvailabilityService)
	userAvailabilityHandler := NewUserAvailabilityHandler(mockUserAvailService)

	t.Run("missing event_id in URL, should return bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/1/availability/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		w := httptest.NewRecorder()
		userAvailabilityHandler.GetUserAvailability(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "event_id is required")
	})

	t.Run("missing user_id in URL, should return bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/1/availability/1", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		w := httptest.NewRecorder()
		userAvailabilityHandler.GetUserAvailability(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user_id is required")
	})

	t.Run("service error, should return internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/1/availability/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1", "event_id": "1"})
		w := httptest.NewRecorder()
		mockUserAvailService.On("GetUserAvailability", req.Context(), int64(1), int64(1)).Return(nil, assert.AnError).Once()

		userAvailabilityHandler.GetUserAvailability(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), assert.AnError.Error())
		mockUserAvailService.AssertExpectations(t)
	})

	t.Run("successful retrieval of user availability", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/1/availability/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1", "event_id": "1"})
		w := httptest.NewRecorder()

		mockResponse := []model.EventSlot{
			{ID: 1, StartTime: time.Date(2025, 07, 13, 10, 0, 0, 0, time.UTC), EndTime: time.Date(2025, 07, 13, 11, 0, 0, 0, time.UTC)},
			{ID: 2, StartTime: time.Date(2025, 07, 13, 12, 0, 0, 0, time.UTC), EndTime: time.Date(2025, 07, 13, 13, 0, 0, 0, time.UTC)},
		}
		mockUserAvailService.On("GetUserAvailability", req.Context(), int64(1), int64(1)).Return(mockResponse, nil).Once()
		userAvailabilityHandler.GetUserAvailability(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockUserAvailService.AssertExpectations(t)
	})

}

func TestDeleteUserAvailability(t *testing.T) {
	mockUserAvailService := new(mockService.MockUserAvailabilityService)
	userAvailabilityHandler := NewUserAvailabilityHandler(mockUserAvailService)

	t.Run("missing event_id in URL, should return bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/1/availability/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		w := httptest.NewRecorder()
		userAvailabilityHandler.DeleteUserAvailability(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "event_id is required")
	})

	t.Run("missing user_id in URL, should return bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/1/availability/1", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		w := httptest.NewRecorder()
		userAvailabilityHandler.DeleteUserAvailability(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user_id is required")
	})

	t.Run("service error, should return internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/1/availability/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1", "event_id": "1"})
		w := httptest.NewRecorder()
		mockUserAvailService.On("DeleteUserAvailability", req.Context(), int64(1), int64(1)).Return(assert.AnError).Once()

		userAvailabilityHandler.DeleteUserAvailability(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), assert.AnError.Error())
		mockUserAvailService.AssertExpectations(t)
	})

	t.Run("successful deletion of user availability", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/1/availability/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1", "event_id": "1"})
		w := httptest.NewRecorder()

		mockUserAvailService.On("DeleteUserAvailability", req.Context(), int64(1), int64(1)).Return(nil).
			Once()
		userAvailabilityHandler.DeleteUserAvailability(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockUserAvailService.AssertExpectations(t)
		assert.Contains(t, w.Body.String(), `{"message":"User availability deleted successfully"}`)
	})

}

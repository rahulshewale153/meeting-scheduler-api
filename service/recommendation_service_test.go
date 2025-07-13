package service

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	mock_repository "github.com/rahulshewale153/meeting-scheduler-api/mock/repository"
	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/assert"
)

func TestGetRecommendedSlots(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	mockUserAvailRepo := new(mock_repository.MockUserAvailabilityRepository)
	mockEventRepo := new(mock_repository.MockEventRepository)
	recommendationService := NewRecommendationService(mockEventRepo, mockUserAvailRepo)
	ctx := context.Background()
	eventID := int64(1)

	t.Run("Function must return an error when the get event operation fails", func(t *testing.T) {
		mockEventRepo.On("GetEvent", ctx, eventID).Return(model.Event{}, assert.AnError).Once()
		_, err := recommendationService.GetRecommendedSlots(ctx, eventID)
		assert.Error(t, err)
		mockEventRepo.AssertExpectations(t)
	})

	t.Run("Function must return an error when the get event slot operation fails", func(t *testing.T) {
		mockEventRepo.On("GetEvent", ctx, eventID).Return(model.Event{}, nil).Once()
		mockEventRepo.On("GetEventSlots", ctx, eventID).Return([]model.EventSlot{}, assert.AnError).Once()

		_, err := recommendationService.GetRecommendedSlots(ctx, eventID)
		assert.Error(t, err)
		mockEventRepo.AssertExpectations(t)
	})

	t.Run("Function must return an error when the get user availability operation fails", func(t *testing.T) {
		eventUserMap := make(map[int64][]model.EventSlot)
		mockEventRepo.On("GetEvent", ctx, eventID).Return(model.Event{}, nil).Once()
		mockEventRepo.On("GetEventSlots", ctx, eventID).Return([]model.EventSlot{{StartTime: time.Date(2025, 07, 13, 10, 0, 0, 0, time.UTC), EndTime: time.Date(2025, 07, 13, 11, 0, 0, 0, time.UTC)}}, nil).Once()
		mockUserAvailRepo.On("GetAllEventUsers", ctx, eventID).Return(eventUserMap, assert.AnError).Once()

		_, err := recommendationService.GetRecommendedSlots(ctx, eventID)
		assert.Error(t, err)
		mockUserAvailRepo.AssertExpectations(t)
	})

	t.Run("Function must return recommended slots when all operations are successful", func(t *testing.T) {
		eventUserMap := map[int64][]model.EventSlot{
			1: {
				{StartTime: time.Date(2025, 07, 13, 10, 0, 0, 0, time.UTC), EndTime: time.Date(2025, 07, 13, 11, 0, 0, 0, time.UTC)},
			},
		}
		mockEventRepo.On("GetEvent", ctx, eventID).Return(model.Event{OrganizerID: 1, DurationMinutes: 30}, nil).Once()
		mockEventRepo.On("GetEventSlots", ctx, eventID).Return([]model.EventSlot{{StartTime: time.Date(2025, 07, 13, 10, 0, 0, 0, time.UTC), EndTime: time.Date(2025, 07, 13, 11, 0, 0, 0, time.UTC)}}, nil).Once()
		mockUserAvailRepo.On("GetAllEventUsers", ctx, eventID).Return(eventUserMap, nil).Once()

		recommendedSlots, err := recommendationService.GetRecommendedSlots(ctx, eventID)
		assert.NoError(t, err)
		assert.NotEmpty(t, recommendedSlots)
		mockEventRepo.AssertExpectations(t)
		mockUserAvailRepo.AssertExpectations(t)
	})

}

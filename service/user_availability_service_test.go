package service

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	mock_repository "github.com/rahulshewale153/meeting-scheduler-api/mock/repository"
	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestInsertUserAvailability(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// Begin a mock transaction
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	mockUserAvailRepo := new(mock_repository.MockUserAvailabilityRepository)
	mockTransactionManager := new(mock_repository.MockTransactionManager)
	userAvailabilityService := NewUserAvailabilityService(mockTransactionManager, mockUserAvailRepo)
	ctx := context.Background()
	userAvailability := model.UserAvailability{
		UserID:  1,
		EventID: 1,
		Availability: []model.EventSlot{
			{
				StartTime: time.Date(2025, 07, 13, 10, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2025, 07, 13, 11, 0, 0, 0, time.UTC),
			},
		},
	}

	t.Run("Function must return an error when the transaction cannot be started", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, assert.AnError).Once()
		err := userAvailabilityService.InsertUserAvailability(ctx, userAvailability)
		assert.Error(t, err)
		mockTransactionManager.AssertExpectations(t)
	})

	t.Run("Function must return an error when the write operation fails", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
		mockUserAvailRepo.On("InsertUserAvailability", ctx, tx, userAvailability.UserID, userAvailability.EventID, userAvailability.Availability[0].StartTime, userAvailability.Availability[0].EndTime).
			Return(int64(0), assert.AnError).Once()

		err := userAvailabilityService.InsertUserAvailability(ctx, userAvailability)
		assert.Error(t, err)
		mockUserAvailRepo.AssertExpectations(t)
		mockTransactionManager.AssertExpectations(t)
		mock.ExpectRollback()
	})

	t.Run("Function must return nil when the insert operation is successful", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
		mockUserAvailRepo.On("InsertUserAvailability", ctx, tx, userAvailability.UserID, userAvailability.EventID, userAvailability.Availability[0].StartTime, userAvailability.Availability[0].EndTime).
			Return(int64(1), nil).Once()

		err := userAvailabilityService.InsertUserAvailability(ctx, userAvailability)
		assert.NoError(t, err)
		mockUserAvailRepo.AssertExpectations(t)
		mockTransactionManager.AssertExpectations(t)
		mock.ExpectCommit()

	})

}

func TestUpdateUserAvailability(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// Begin a mock transaction
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	mockUserAvailRepo := new(mock_repository.MockUserAvailabilityRepository)
	mockTransactionManager := new(mock_repository.MockTransactionManager)
	userAvailabilityService := NewUserAvailabilityService(mockTransactionManager, mockUserAvailRepo)
	ctx := context.Background()
	userAvailability := model.UserAvailability{
		UserID:  1,
		EventID: 1,
		Availability: []model.EventSlot{
			{
				StartTime: time.Date(2025, 07, 12, 12, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2025, 07, 12, 13, 0, 0, 0, time.UTC),
			},
		},
	}

	t.Run("Function must return an error when the transaction cannot be started", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, assert.AnError).Once()
		err := userAvailabilityService.UpdateUserAvailability(ctx, userAvailability)
		assert.Error(t, err)
		mockTransactionManager.AssertExpectations(t)
	})

	t.Run("Function must return an error when GetUserAvailability operation fails", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
		mockUserAvailRepo.On("GetUserAvailability", ctx, userAvailability.EventID, userAvailability.UserID).
			Return(nil, assert.AnError).Once()
		err := userAvailabilityService.UpdateUserAvailability(ctx, userAvailability)
		assert.Error(t, err)
		mockUserAvailRepo.AssertExpectations(t)
		mockTransactionManager.AssertExpectations(t)
		mock.ExpectRollback()
	})

	t.Run("Function must return an error when the GetUserAvailability operation is successful", func(t *testing.T) {
		t.Run("Function must return an error when the update operation fails", func(t *testing.T) {
			mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
			mockUserAvailRepo.On("GetUserAvailability", ctx, userAvailability.EventID, userAvailability.UserID).Return([]model.EventSlot{model.EventSlot{ID: 1, StartTime: time.Date(2025, 07, 12, 10, 0, 0, 0, time.UTC), EndTime: time.Date(2025, 07, 12, 11, 0, 0, 0, time.UTC)}}, nil).Once()
			mockUserAvailRepo.On("InsertUserAvailability", ctx, tx, userAvailability.UserID, userAvailability.EventID, testifyMock.Anything, testifyMock.Anything).
				Return(int64(0), assert.AnError).Once()

			err := userAvailabilityService.UpdateUserAvailability(ctx, userAvailability)
			assert.Error(t, err)
			mockUserAvailRepo.AssertExpectations(t)
			mockTransactionManager.AssertExpectations(t)
			mock.ExpectRollback()
		})

		t.Run("Function must return nil when the update operation is successful", func(t *testing.T) {
			mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
			mockUserAvailRepo.On("GetUserAvailability", ctx, userAvailability.EventID, userAvailability.UserID).Return([]model.EventSlot{model.EventSlot{ID: 1, StartTime: time.Date(2025, 07, 12, 10, 0, 0, 0, time.UTC), EndTime: time.Date(2025, 07, 12, 11, 0, 0, 0, time.UTC)}}, nil).Once()
			mockUserAvailRepo.On("InsertUserAvailability", ctx, tx, userAvailability.UserID, userAvailability.EventID, testifyMock.Anything, testifyMock.Anything).
				Return(int64(1), nil).Once()
			mockUserAvailRepo.On("DeleteUserAvailability", ctx, tx, userAvailability.UserID, int64(1)).
				Return(nil).Once()
			err := userAvailabilityService.UpdateUserAvailability(ctx, userAvailability)
			assert.NoError(t, err)
			mockUserAvailRepo.AssertExpectations(t)
			mockTransactionManager.AssertExpectations(t)
			mock.ExpectCommit()
		})
	})
}

func TestDeleteUserAvailability(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// Begin a mock transaction
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	mockUserAvailRepo := new(mock_repository.MockUserAvailabilityRepository)
	mockTransactionManager := new(mock_repository.MockTransactionManager)
	userAvailabilityService := NewUserAvailabilityService(mockTransactionManager, mockUserAvailRepo)
	ctx := context.Background()
	userID := int64(1)
	eventID := int64(1)

	t.Run("Function must return an error when the transaction cannot be started", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, assert.AnError).Once()
		err := userAvailabilityService.DeleteUserAvailability(ctx, userID, eventID)
		assert.Error(t, err)
		mockTransactionManager.AssertExpectations(t)
	})

	t.Run("Function must return an error when the delete operation fails", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
		mockUserAvailRepo.On("DeleteUserAvailability", ctx, tx, userID, eventID).
			Return(assert.AnError).Once()

		err := userAvailabilityService.DeleteUserAvailability(ctx, userID, eventID)
		assert.Error(t, err)
		mockUserAvailRepo.AssertExpectations(t)
		mockTransactionManager.AssertExpectations(t)
		mock.ExpectRollback()
	})

	t.Run("Function must return nil when the delete operation is successful", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
		mockUserAvailRepo.On("DeleteUserAvailability", ctx, tx, userID, eventID).
			Return(nil).Once()

		err := userAvailabilityService.DeleteUserAvailability(ctx, userID, eventID)
		assert.NoError(t, err)
		mockUserAvailRepo.AssertExpectations(t)
		mockTransactionManager.AssertExpectations(t)
		mock.ExpectCommit()
	})

}

func TestGetUserAvailability(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	mockUserAvailRepo := new(mock_repository.MockUserAvailabilityRepository)
	userAvailabilityService := NewUserAvailabilityService(nil, mockUserAvailRepo)
	ctx := context.Background()
	eventID := int64(1)
	userID := int64(1)

	t.Run("Function must return an error when the read operation fails", func(t *testing.T) {
		mockUserAvailRepo.On("GetUserAvailability", ctx, eventID, userID).
			Return(nil, assert.AnError).Once()

		_, err := userAvailabilityService.GetUserAvailability(ctx, eventID, userID)
		assert.Error(t, err)
		mockUserAvailRepo.AssertExpectations(t)
	})

	t.Run("Function must return nil when the read operation is successful", func(t *testing.T) {
		expectedSlots := []model.EventSlot{
			{StartTime: time.Date(2025, 07, 13, 10, 0, 0, 0, time.UTC), EndTime: time.Date(2025, 07, 13, 11, 0, 0, 0, time.UTC)},
		}
		mockUserAvailRepo.On("GetUserAvailability", ctx, eventID, userID).
			Return(expectedSlots, nil).Once()

		slots, err := userAvailabilityService.GetUserAvailability(ctx, eventID, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedSlots, slots)
		mockUserAvailRepo.AssertExpectations(t)
	})
}

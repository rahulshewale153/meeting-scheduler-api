package service

import (
	"context"
	"testing"

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
				StartTime: "2023-10-01 10:00:00",
				EndTime:   "2023-10-01 11:00:00",
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
				StartTime: "2025-07-12T10:00:00Z",
				EndTime:   "2025-07-12T11:00:00Z",
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
			mockUserAvailRepo.On("GetUserAvailability", ctx, userAvailability.EventID, userAvailability.UserID).Return([]model.EventSlot{model.EventSlot{ID: 1, StartTime: "2023-10-01 10:00:00", EndTime: "2023-10-01 11:00:00"}}, nil).Once()
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
			mockUserAvailRepo.On("GetUserAvailability", ctx, userAvailability.EventID, userAvailability.UserID).Return([]model.EventSlot{model.EventSlot{ID: 1, StartTime: "2023-10-01 10:00:00", EndTime: "2023-10-01 11:00:00"}}, nil).Once()
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

}

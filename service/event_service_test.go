package service

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	mock_repository "github.com/rahulshewale153/meeting-scheduler-api/mock/repository"
	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/assert"
)

func TestInsertEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// Begin a mock transaction
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	mockEventRepo := new(mock_repository.MockEventRepository)
	mockTransactionManager := new(mock_repository.MockTransactionManager)
	service := NewEventService(mockTransactionManager, mockEventRepo)
	ctx := context.Background()
	createEventReq := model.EventRequest{
		Event: model.Event{
			Title:           "Test Event",
			OrganizerID:     1,
			DurationMinutes: 60,
		},
		ProposedSlots: []model.EventSlot{
			{
				StartTime: "2025-07-12T10:00:00Z",
				EndTime:   "2025-07-12T11:00:00Z",
			},
		},
	}

	t.Run("Function must return an error when the transaction cannot be started", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, assert.AnError).Once()
		_, err := service.InsertEvent(ctx, createEventReq)
		assert.Error(t, err)
		mockTransactionManager.AssertExpectations(t)
	})

	t.Run("Function must return an error when the insert operation fails", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
		mockEventRepo.On("InsertEvent", ctx, tx, createEventReq.Event).
			Return(int64(0), assert.AnError).Once()

		_, err := service.InsertEvent(ctx, createEventReq)
		assert.Error(t, err)
		mockTransactionManager.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
		mock.ExpectRollback()
	})

	t.Run("Function must return the event_id when the insert operation is successful", func(t *testing.T) {

		t.Run("Function must return an error when the slot start time conversion to UTC fails", func(t *testing.T) {
			mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
			mockEventRepo.On("InsertEvent", ctx, tx, createEventReq.Event).
				Return(int64(1), nil).Once()
			createEventReq.ProposedSlots[0].StartTime = "invalid-time"
			_, err := service.InsertEvent(ctx, createEventReq)
			assert.Error(t, err)
			mockTransactionManager.AssertExpectations(t)
			mockEventRepo.AssertExpectations(t)
			mock.ExpectRollback()
		})

		t.Run("Function must return an error when the slot end time conversion to UTC fails", func(t *testing.T) {
			mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
			mockEventRepo.On("InsertEvent", ctx, tx, createEventReq.Event).
				Return(int64(1), nil).Once()
			createEventReq.ProposedSlots[0].StartTime = "2025-07-12T10:00:00Z"
			createEventReq.ProposedSlots[0].EndTime = "invalid-time"
			_, err := service.InsertEvent(ctx, createEventReq)
			assert.Error(t, err)
			mockTransactionManager.AssertExpectations(t)
			mockEventRepo.AssertExpectations(t)
			mock.ExpectRollback()
		})

		t.Run("Function must return an error when the insert event slots operation fails", func(t *testing.T) {
			createEventReq.ProposedSlots[0].StartTime = "2025-07-12T10:00:00Z"
			createEventReq.ProposedSlots[0].EndTime = "2025-07-12T11:00:00Z"
			mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
			mockEventRepo.On("InsertEvent", ctx, tx, createEventReq.Event).
				Return(int64(1), nil).Once()
			mockEventRepo.On("InsertEventSlots", ctx, tx, int64(1), createEventReq.ProposedSlots[0]).
				Return(assert.AnError).Once()
			_, err := service.InsertEvent(ctx, createEventReq)
			assert.Error(t, err)
			mockTransactionManager.AssertExpectations(t)
			mockEventRepo.AssertExpectations(t)
			mock.ExpectRollback()
		})

		t.Run("Function must return nil when the insert operation is successful", func(t *testing.T) {
			createEventReq.ProposedSlots[0].StartTime = "2025-07-12T10:00:00Z"
			createEventReq.ProposedSlots[0].EndTime = "2025-07-12T11:00:00Z"
			mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
			mockEventRepo.On("InsertEvent", ctx, tx, createEventReq.Event).
				Return(int64(1), nil).Once()
			mockEventRepo.On("InsertEventSlots", ctx, tx, int64(1), createEventReq.ProposedSlots[0]).
				Return(nil).Once()

			eventID, err := service.InsertEvent(ctx, createEventReq)
			assert.NoError(t, err)
			assert.Equal(t, int64(1), eventID)
			mockTransactionManager.AssertExpectations(t)
			mockEventRepo.AssertExpectations(t)
			mock.ExpectCommit()
		})

	})

}

func TestUpdateEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// Begin a mock transaction
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	mockEventRepo := new(mock_repository.MockEventRepository)
	mockTransactionManager := new(mock_repository.MockTransactionManager)
	service := NewEventService(mockTransactionManager, mockEventRepo)
	ctx := context.Background()
	updateEventReq := model.EventRequest{
		Event: model.Event{
			ID:              1,
			Title:           "Updated Event",
			OrganizerID:     2,
			DurationMinutes: 90,
		},
		ProposedSlots: []model.EventSlot{
			{
				StartTime: "2025-07-12T12:00:00Z",
				EndTime:   "2025-07-12T13:00:00Z",
			},
		},
	}

	t.Run("Function must return an error when the transaction cannot be started", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, assert.AnError).Once()
		err := service.UpdateEvent(ctx, updateEventReq)
		assert.Error(t, err)
		mockTransactionManager.AssertExpectations(t)
	})

	t.Run("Function must return an error when the update operation fails", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
		mockEventRepo.On("UpdateEvent", ctx, tx, updateEventReq.Event).
			Return(assert.AnError).Once()

		err := service.UpdateEvent(ctx, updateEventReq)
		assert.Error(t, err)
		mockTransactionManager.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
		mock.ExpectRollback()
	})

	t.Run("Function must return nil when the update operation is successful", func(t *testing.T) {
		t.Run("Function must return an error when the get slot operation fails", func(t *testing.T) {
			mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
			mockEventRepo.On("UpdateEvent", ctx, tx, updateEventReq.Event).
				Return(nil).Once()
			mockEventRepo.On("GetEventSlots", ctx, updateEventReq.Event.ID).
				Return([]model.EventSlot{}, assert.AnError).Once()

			err := service.UpdateEvent(ctx, updateEventReq)
			assert.Error(t, err)
			mockTransactionManager.AssertExpectations(t)
			mockEventRepo.AssertExpectations(t)
			mock.ExpectRollback()
		})
		t.Run("Function must return an error when the insert slot operation fails", func(t *testing.T) {
			mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
			mockEventRepo.On("UpdateEvent", ctx, tx, updateEventReq.Event).
				Return(nil).Once()
			mockEventRepo.On("GetEventSlots", ctx, updateEventReq.Event.ID).
				Return([]model.EventSlot{}, nil).Once()
			mockEventRepo.On("InsertEventSlots", ctx, tx, updateEventReq.Event.ID, model.EventSlot{
				StartTime: "2025-07-12T12:00:00Z",
				EndTime:   "2025-07-12T13:00:00Z",
			}).Return(assert.AnError).Once()

			err := service.UpdateEvent(ctx, updateEventReq)
			assert.Error(t, err)
			mockTransactionManager.AssertExpectations(t)
			mockEventRepo.AssertExpectations(t)
			mock.ExpectRollback()
		})

		t.Run("Function must return nil when the update operation is successful", func(t *testing.T) {
			mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
			mockEventRepo.On("UpdateEvent", ctx, tx, updateEventReq.Event).
				Return(nil).Once()
			mockEventRepo.On("GetEventSlots", ctx, updateEventReq.Event.ID).
				Return([]model.EventSlot{}, nil).Once()
			mockEventRepo.On("InsertEventSlots", ctx, tx, updateEventReq.Event.ID, model.EventSlot{
				StartTime: "2025-07-12T12:00:00Z",
				EndTime:   "2025-07-12T13:00:00Z",
			}).Return(nil).Once()

			err := service.UpdateEvent(ctx, updateEventReq)
			assert.NoError(t, err)
			mockTransactionManager.AssertExpectations(t)
			mockEventRepo.AssertExpectations(t)
			mock.ExpectCommit()
		})

	})
}

func TestDeleteEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// Begin a mock transaction
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	mockEventRepo := new(mock_repository.MockEventRepository)
	mockTransactionManager := new(mock_repository.MockTransactionManager)
	service := NewEventService(mockTransactionManager, mockEventRepo)
	ctx := context.Background()
	eventID := int64(1)

	t.Run("Function must return an error when the transaction cannot be started", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, assert.AnError).Once()
		err := service.DeleteEvent(ctx, eventID)
		assert.Error(t, err)
		mockTransactionManager.AssertExpectations(t)
	})

	t.Run("Function must return an error when the delete event slot operation fails", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
		mockEventRepo.On("DeleteEventSlots", ctx, tx, eventID).
			Return(assert.AnError).Once()

		err := service.DeleteEvent(ctx, eventID)
		assert.Error(t, err)
		mockTransactionManager.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
		mock.ExpectRollback()
	})

	t.Run("Function must return an error when the delete operation fails", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
		mockEventRepo.On("DeleteEventSlots", ctx, tx, eventID).
			Return(nil).Once()

		mockEventRepo.On("DeleteEvent", ctx, tx, eventID).
			Return(nil).Once()

		err := service.DeleteEvent(ctx, eventID)
		assert.NoError(t, err)
		mockTransactionManager.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
		mock.ExpectRollback()
	})

	t.Run("Function must return nil when the delete operation is successful", func(t *testing.T) {
		mockTransactionManager.On("BeginTransaction", ctx).Return(tx, nil).Once()
		mockEventRepo.On("DeleteEventSlots", ctx, tx, eventID).
			Return(nil).Once()

		mockEventRepo.On("DeleteEvent", ctx, tx, eventID).
			Return(nil).Once()

		err := service.DeleteEvent(ctx, eventID)
		assert.NoError(t, err)
		mockTransactionManager.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
		mock.ExpectCommit()
	})

}

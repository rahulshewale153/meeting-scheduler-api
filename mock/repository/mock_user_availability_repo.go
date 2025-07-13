package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/mock"
)

type MockUserAvailabilityRepository struct {
	mock.Mock
}

func (m *MockUserAvailabilityRepository) InsertUserAvailability(ctx context.Context, tx *sql.Tx, userID int64, eventID int64, startTime time.Time, endTime time.Time) (int64, error) {
	args := m.Called(ctx, tx, userID, eventID, startTime, endTime)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserAvailabilityRepository) GetAllEventUsers(ctx context.Context, eventID int64) (map[int64][]model.EventSlot, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).(map[int64][]model.EventSlot), args.Error(1)
}

func (m *MockUserAvailabilityRepository) DeleteUserAvailability(ctx context.Context, tx *sql.Tx, userID int64, availabilityID int64) error {

	args := m.Called(ctx, tx, userID, availabilityID)
	return args.Error(0)
}

func (m *MockUserAvailabilityRepository) GetUserAvailability(ctx context.Context, eventID int64, userID int64) ([]model.EventSlot, error) {
	args := m.Called(ctx, eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.EventSlot), args.Error(1)
}

func (m *MockUserAvailabilityRepository) GetEventUsers(ctx context.Context, eventID int64) (map[int64][]model.EventSlot, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[int64][]model.EventSlot), args.Error(1)
}

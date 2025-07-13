package service

import (
	"context"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/mock"
)

type MockUserAvailabilityService struct {
	mock.Mock
}

func (m *MockUserAvailabilityService) InsertUserAvailability(ctx context.Context, userAvailability model.UserAvailability) error {
	args := m.Called(ctx, userAvailability)
	return args.Error(0)
}

func (m *MockUserAvailabilityService) UpdateUserAvailability(ctx context.Context, userAvailability model.UserAvailability) error {
	args := m.Called(ctx, userAvailability)
	return args.Error(0)
}

func (m *MockUserAvailabilityService) DeleteUserAvailability(ctx context.Context, userID int64, availabilityID int64) error {
	args := m.Called(ctx, userID, availabilityID)
	return args.Error(0)
}

func (m *MockUserAvailabilityService) GetUserAvailability(ctx context.Context, eventID int64, userID int64) ([]model.EventSlot, error) {
	args := m.Called(ctx, eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.EventSlot), args.Error(1)
}

func (m *MockUserAvailabilityService) GetEventUsers(ctx context.Context, eventID int64) (map[int64][]model.EventSlot, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[int64][]model.EventSlot), args.Error(1)
}

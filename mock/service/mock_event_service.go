package service

import (
	"context"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/mock"
)

type MockEventService struct {
	mock.Mock
}

func (m *MockEventService) InsertEvent(ctx context.Context, createEventReq model.EventRequest) (int64, error) {
	args := m.Called(ctx, createEventReq)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEventService) UpdateEvent(ctx context.Context, updateEventReq model.EventRequest) error {
	args := m.Called(ctx, updateEventReq)
	return args.Error(0)
}

func (m *MockEventService) DeleteEvent(ctx context.Context, eventID int64) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

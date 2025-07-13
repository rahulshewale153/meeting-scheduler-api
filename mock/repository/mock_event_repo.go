package repository

import (
	"context"
	"database/sql"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/mock"
)

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) InsertEvent(ctx context.Context, tx *sql.Tx, createEventReq model.Event) (int64, error) {
	args := m.Called(ctx, tx, createEventReq)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEventRepository) UpdateEvent(ctx context.Context, tx *sql.Tx, updateEventReq model.Event) error {
	args := m.Called(ctx, tx, updateEventReq)
	return args.Error(0)
}

func (m *MockEventRepository) DeleteEvent(ctx context.Context, tx *sql.Tx, eventID int64) error {
	args := m.Called(ctx, tx, eventID)
	return args.Error(0)
}

func (m *MockEventRepository) InsertEventSlots(ctx context.Context, tx *sql.Tx, eventID int64, slot model.EventSlot) error {
	args := m.Called(ctx, tx, eventID, slot)
	return args.Error(0)
}

func (m *MockEventRepository) DeleteEventSlots(ctx context.Context, tx *sql.Tx, slotID int64) error {
	args := m.Called(ctx, tx, slotID)
	return args.Error(0)
}

func (m *MockEventRepository) GetEventSlots(ctx context.Context, eventID int64) ([]model.EventSlot, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).([]model.EventSlot), args.Error(1)
}

func (m *MockEventRepository) GetEvent(ctx context.Context, eventID int64) (model.Event, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).(model.Event), args.Error(1)
}

package service

import (
	"context"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/mock"
)

type MockRecommendationService struct {
	mock.Mock
}

func (m *MockRecommendationService) GetRecommendedSlots(ctx context.Context, eventID int64) ([]model.SlotRecommendation, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).([]model.SlotRecommendation), args.Error(1)
}

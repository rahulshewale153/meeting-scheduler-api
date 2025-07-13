package service

import (
	"context"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
)

type EventServiceI interface {
	InsertEvent(ctx context.Context, createEventReq model.EventRequest) (int64, error)
	UpdateEvent(ctx context.Context, updateEventReq model.EventRequest) error
	DeleteEvent(ctx context.Context, eventID int64) error
}

type UserAvailabilityServiceI interface {
	InsertUserAvailability(ctx context.Context, userAvailability model.UserAvailability) error
	UpdateUserAvailability(ctx context.Context, userAvailability model.UserAvailability) error
	DeleteUserAvailability(ctx context.Context, userID int64, eventID int64) error
	GetUserAvailability(ctx context.Context, eventID int64, userID int64) ([]model.EventSlot, error)
}

type RecommendationServiceI interface {
	GetRecommendedSlots(ctx context.Context, eventID int64) ([]model.SlotRecommendation, error)
}

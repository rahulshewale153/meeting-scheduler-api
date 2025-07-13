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

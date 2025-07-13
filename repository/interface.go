package repository

import (
	"context"
	"database/sql"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
)

type TransactionManagerI interface {
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
}

type EventRepositoryI interface {
	InsertEvent(ctx context.Context, tx *sql.Tx, createEventReq model.Event) (int64, error)
	UpdateEvent(ctx context.Context, tx *sql.Tx, updateEventReq model.Event) error
	DeleteEvent(ctx context.Context, tx *sql.Tx, eventID int64) error
	InsertEventSlots(ctx context.Context, tx *sql.Tx, eventID int64, slot model.EventSlot) error
	DeleteEventSlots(ctx context.Context, tx *sql.Tx, slotID int64) error
	GetEventSlots(ctx context.Context, eventID int64) ([]model.EventSlot, error)
}

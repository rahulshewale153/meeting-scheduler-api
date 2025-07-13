package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

	repository := NewEventRepository(db)
	ctx := context.Background()
	createEventReq := model.Event{
		Title:           "Test Event",
		OrganizerID:     1,
		DurationMinutes: 60,
	}
	query := `INSERT INTO event_detail (title, organizer_id, duration_minutes) VALUES (?, ?, ?)`
	t.Run("Function must return an error when the write operation fails", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(createEventReq.Title, createEventReq.OrganizerID, createEventReq.DurationMinutes).
			WillReturnError(assert.AnError)

		_, err := repository.InsertEvent(ctx, tx, createEventReq)
		assert.Error(t, err)
	})

	t.Run("Function must return the event_id when the insert operation is successful", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(createEventReq.Title, createEventReq.OrganizerID, createEventReq.DurationMinutes).
			WillReturnResult(sqlmock.NewResult(1, 1))

		eventID, err := repository.InsertEvent(ctx, tx, createEventReq)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), eventID)
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

	repository := NewEventRepository(db)
	ctx := context.Background()
	updateEventReq := model.Event{
		ID:              1,
		Title:           "Updated Event",
		OrganizerID:     1,
		DurationMinutes: 90,
	}

	query := `Update event_detail SET title = ?, organizer_id = ?, duration_minutes = ? WHERE id = ?`
	t.Run("Function must return an error when the write operation fails", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(updateEventReq.Title, updateEventReq.OrganizerID, updateEventReq.DurationMinutes, updateEventReq.ID).
			WillReturnError(assert.AnError)

		err := repository.UpdateEvent(ctx, tx, updateEventReq)
		assert.Error(t, err)
	})

	t.Run("Function must return nil when the update operation is successful", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(updateEventReq.Title, updateEventReq.OrganizerID, updateEventReq.DurationMinutes, updateEventReq.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repository.UpdateEvent(ctx, tx, updateEventReq)
		assert.NoError(t, err)
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

	repository := NewEventRepository(db)
	ctx := context.Background()
	eventID := int64(1)

	query := `DELETE FROM event_detail WHERE id = ?`
	t.Run("Function must return an error when the delete operation fails", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(eventID).
			WillReturnError(assert.AnError)

		err := repository.DeleteEvent(ctx, tx, eventID)
		assert.Error(t, err)
	})

	t.Run("Function must return nil when the delete operation is successful", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(eventID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repository.DeleteEvent(ctx, tx, eventID)
		assert.NoError(t, err)
	})

}

func TestInsertEventSlots(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// Begin a mock transaction
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	repository := NewEventRepository(db)
	ctx := context.Background()
	eventID := int64(1)
	slot := model.EventSlot{
		StartTime: "2023-10-01T10:00:00Z",
		EndTime:   "2023-10-01T11:00:00Z",
	}

	query := `INSERT INTO event_slot (event_id, start_time, end_time) VALUES (?, ?, ?)`
	t.Run("Function must return an error when the write operation fails", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(eventID, slot.StartTime, slot.EndTime).
			WillReturnError(assert.AnError)

		err := repository.InsertEventSlots(ctx, tx, eventID, slot)
		assert.Error(t, err)
	})

	t.Run("Function must return nil when the insert operation is successful", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(eventID, slot.StartTime, slot.EndTime).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repository.InsertEventSlots(ctx, tx, eventID, slot)
		assert.NoError(t, err)
	})

}
func TestDeleteEventSlots(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// Begin a mock transaction
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	repository := NewEventRepository(db)
	ctx := context.Background()
	slotID := int64(1)

	query := `DELETE FROM event_slot WHERE id = ?`
	t.Run("Function must return an error when the delete operation fails", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(slotID).
			WillReturnError(assert.AnError)

		err := repository.DeleteEventSlots(ctx, tx, slotID)
		assert.Error(t, err)
	})

	t.Run("Function must return nil when the delete operation is successful", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(slotID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repository.DeleteEventSlots(ctx, tx, slotID)
		assert.NoError(t, err)
	})

}

func TestGetEventSlots(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repository := NewEventRepository(db)
	ctx := context.Background()
	eventID := int64(1)

	query := `SELECT id, start_time, end_time FROM event_slot WHERE event_id = ?`
	t.Run("Function must return an error when the read operation fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(eventID).
			WillReturnError(assert.AnError)

		_, err := repository.GetEventSlots(ctx, eventID)
		assert.Error(t, err)
	})

	t.Run("Function must return an error when scanning the rows fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(eventID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "start_time", "end_time"}).
				AddRow(nil, "2023-10-01T10:00:00Z", "2023-10-01T11:00:00Z"))
		_, err := repository.GetEventSlots(ctx, eventID)
		assert.Error(t, err)
	})

	t.Run("Function must return an empty slice when no slots are found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(eventID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "start_time", "end_time"}))
		slots, err := repository.GetEventSlots(ctx, eventID)
		assert.NoError(t, err)
		assert.Empty(t, slots)
	})

	t.Run("Function must return the event slots when the read operation is successful", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "start_time", "end_time"}).
			AddRow(1, "2023-10-01T10:00:00Z", "2023-10-01T11:00:00Z").
			AddRow(2, "2023-10-01T12:00:00Z", "2023-10-01T13:00:00Z")

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(eventID).
			WillReturnRows(rows)

		slots, err := repository.GetEventSlots(ctx, eventID)
		assert.NoError(t, err)
		assert.Len(t, slots, 2)
		assert.Equal(t, int64(1), slots[0].ID)
		assert.Equal(t, int64(2), slots[1].ID)
	})

}

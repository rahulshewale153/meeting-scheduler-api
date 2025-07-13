package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/stretchr/testify/assert"
)

func TestInsertUserAvailability(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// Begin a mock transaction
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	repository := NewUserAvailabilityRepository(db)
	ctx := context.Background()

	createUserAvailabilityReq := model.UserAvailability{
		UserID:  1,
		EventID: 1,
		Availability: []model.EventSlot{
			{
				StartTime: "2023-10-01 10:00:00",
				EndTime:   "2023-10-01 11:00:00",
			},
		},
	}

	query := `INSERT INTO user_availability (event_id, user_id, start_time, end_time) VALUES (?, ?, ?, ?)`
	t.Run("Function must return an error when the write operation fails", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(createUserAvailabilityReq.EventID, createUserAvailabilityReq.UserID, createUserAvailabilityReq.Availability[0].StartTime, createUserAvailabilityReq.Availability[0].EndTime).
			WillReturnError(assert.AnError)

		_, err := repository.InsertUserAvailability(ctx, tx, createUserAvailabilityReq.UserID, createUserAvailabilityReq.EventID, createUserAvailabilityReq.Availability[0].StartTime, createUserAvailabilityReq.Availability[0].EndTime)
		assert.Error(t, err)
	})

	t.Run("Function must return the last inserted ID when the insert operation is successful", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(createUserAvailabilityReq.EventID, createUserAvailabilityReq.UserID, createUserAvailabilityReq.Availability[0].StartTime, createUserAvailabilityReq.Availability[0].EndTime).
			WillReturnResult(sqlmock.NewResult(1, 1))

		lastInsertID, err := repository.InsertUserAvailability(ctx, tx, createUserAvailabilityReq.UserID, createUserAvailabilityReq.EventID, createUserAvailabilityReq.Availability[0].StartTime, createUserAvailabilityReq.Availability[0].EndTime)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), lastInsertID)
	})

}

func TestGetAllEventUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repository := NewUserAvailabilityRepository(db)
	ctx := context.Background()

	eventID := int64(1)

	query := `SELECT id, user_id, start_time, end_time FROM user_availability WHERE event_id = ?`
	t.Run("Function must return an error when the read operation fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(eventID).
			WillReturnError(assert.AnError)

		_, err := repository.GetAllEventUsers(ctx, eventID)
		assert.Error(t, err)
	})

	t.Run("Function must return a map of user IDs and their availability when the read operation is successful", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "start_time", "end_time"}).
			AddRow(1, 1, "2023-10-01 10:00:00", "2023-10-01 11:00:00").
			AddRow(2, 2, "2023-10-01 12:00:00", "2023-10-01 13:00:00")

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(eventID).
			WillReturnRows(rows)

		eventUsers, err := repository.GetAllEventUsers(ctx, eventID)
		assert.NoError(t, err)
		assert.Len(t, eventUsers, 2)
		assert.Contains(t, eventUsers[1], model.EventSlot{ID: 1, StartTime: "2023-10-01 10:00:00", EndTime: "2023-10-01 11:00:00"})
		assert.Contains(t, eventUsers[2], model.EventSlot{ID: 2, StartTime: "2023-10-01 12:00:00", EndTime: "2023-10-01 13:00:00"})
	})
}

func TestDeleteUserAvailability(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// Begin a mock transaction
	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	repository := NewUserAvailabilityRepository(db)
	ctx := context.Background()

	userID := int64(1)
	availabilityID := int64(1)

	query := `DELETE FROM user_availability WHERE id = ? AND user_id = ?`
	t.Run("Function must return an error when the delete operation fails", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(userID, availabilityID).
			WillReturnError(assert.AnError)

		err := repository.DeleteUserAvailability(ctx, tx, userID, availabilityID)
		assert.Error(t, err)
	})

	t.Run("Function must not return an error when the delete operation is successful", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(userID, availabilityID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repository.DeleteUserAvailability(ctx, tx, userID, availabilityID)
		assert.NoError(t, err)
	})

}

func TestGetUserAvailability(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repository := NewUserAvailabilityRepository(db)
	ctx := context.Background()

	eventID := int64(1)
	userID := int64(1)

	query := `SELECT id, start_time, end_time FROM user_availability WHERE event_id = ? AND user_id = ?`
	t.Run("Function must return an error when the read operation fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(eventID, userID).
			WillReturnError(assert.AnError)

		_, err := repository.GetUserAvailability(ctx, eventID, userID)
		assert.Error(t, err)
	})

	t.Run("Function must return a slice of EventSlot when the read operation is successful", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "start_time", "end_time"}).
			AddRow(1, "2023-10-01 10:00:00", "2023-10-01 11:00:00").
			AddRow(2, "2023-10-01 12:00:00", "2023-10-01 13:00:00")

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(eventID, userID).
			WillReturnRows(rows)

		slots, err := repository.GetUserAvailability(ctx, eventID, userID)
		assert.NoError(t, err)
		assert.Len(t, slots, 2)
		assert.Equal(t, slots[0].StartTime, "2023-10-01 10:00:00")
		assert.Equal(t, slots[0].EndTime, "2023-10-01 11:00:00")
	})
}

package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
)

type userAvailabilityRepository struct {
	dbConn *sql.DB
}

func NewUserAvailabilityRepository(dbConn *sql.DB) UserAvailabilityRepositoryI {
	return &userAvailabilityRepository{dbConn: dbConn}
}

// InsertUserAvailability: inserts a new user availability record into the database.
func (userRepo *userAvailabilityRepository) InsertUserAvailability(ctx context.Context, tx *sql.Tx, userID int64, eventID int64, startTime time.Time, endTime time.Time) (int64, error) {
	query := `INSERT INTO user_availability (event_id, user_id, start_time, end_time) VALUES (?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, query, eventID, userID, startTime, endTime)
	if err != nil {
		log.Printf("Error inserting user availability: %v", err)
		return 0, err
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		return 0, err
	}
	return lastInsertID, nil
}

// GetEventUsers: retrieves the availability of users for a specific event.
func (userRepo *userAvailabilityRepository) GetAllEventUsers(ctx context.Context, eventID int64) (map[int64][]model.EventSlot, error) {
	eventUsers := make(map[int64][]model.EventSlot)
	query := `SELECT id, user_id, start_time, end_time FROM user_availability WHERE event_id = ? order by user_id ASC`
	rows, err := userRepo.dbConn.QueryContext(ctx, query, eventID)
	if err != nil {
		log.Printf("Error retrieving user availability: %v", err)
		return eventUsers, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int64
		var selectedSlot model.EventSlot
		if err := rows.Scan(&selectedSlot.ID, &userID, &selectedSlot.StartTime, &selectedSlot.EndTime); err != nil {
			log.Printf("Error scanning user availability: %v", err)
			return eventUsers, err
		}
		eventUsers[userID] = append(eventUsers[userID], selectedSlot)
	}

	return eventUsers, nil
}

// DeleteUserAvailability: deletes a user availability record from the database.
func (userRepo *userAvailabilityRepository) DeleteUserAvailability(ctx context.Context, tx *sql.Tx, userID int64, eventID int64) error {
	query := `DELETE FROM user_availability WHERE event_id = ? AND user_id = ?`
	_, err := tx.ExecContext(ctx, query, eventID, userID)
	if err != nil {
		log.Printf("Error deleting user availability: %v", err)
		return err
	}
	return nil
}

// GetUserAvailability: retrieves the availability of specific user for a specific event.
func (userRepo *userAvailabilityRepository) GetUserAvailability(ctx context.Context, eventID int64, userID int64) ([]model.EventSlot, error) {
	slots := []model.EventSlot{}
	query := `SELECT id, start_time, end_time FROM user_availability WHERE event_id = ? AND user_id = ?`
	rows, err := userRepo.dbConn.QueryContext(ctx, query, eventID, userID)
	if err != nil {
		log.Printf("Error retrieving user availability: %v", err)
		return slots, err
	}
	defer rows.Close()

	for rows.Next() {
		var slot model.EventSlot
		if err := rows.Scan(&slot.ID, &slot.StartTime, &slot.EndTime); err != nil {
			log.Printf("Error scanning user availability: %v", err)
			return slots, err
		}
		slots = append(slots, slot)
	}

	return slots, nil
}

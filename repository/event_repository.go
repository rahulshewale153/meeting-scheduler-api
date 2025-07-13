package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
)

type eventRepository struct {
	dbConn *sql.DB
}

func NewEventRepository(dbConn *sql.DB) EventRepositoryI {
	return &eventRepository{dbConn: dbConn}
}

// Insert the event
func (eventRepo *eventRepository) InsertEvent(ctx context.Context, tx *sql.Tx, createEventReq model.Event) (int64, error) {
	result, err := tx.ExecContext(ctx, `
		INSERT INTO event_detail (title, organizer_id, duration_minutes) 
		VALUES (?, ?, ?)`, createEventReq.Title, createEventReq.OrganizerID, createEventReq.DurationMinutes)
	if err != nil {
		log.Println("Error inserting event:", err)
		return 0, err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return 0, err
	}

	return eventID, nil
}

// Update the event
func (eventRepo *eventRepository) UpdateEvent(ctx context.Context, tx *sql.Tx, updateEventReq model.Event) error {
	_, err := tx.ExecContext(ctx, `Update event_detail SET title = ?, organizer_id = ?, duration_minutes = ? WHERE id = ?`, updateEventReq.Title, updateEventReq.OrganizerID, updateEventReq.DurationMinutes, updateEventReq.ID)
	if err != nil {
		log.Println("Error updating event:", err)
		return err
	}
	return nil
}

// Delete the event
func (eventRepo *eventRepository) DeleteEvent(ctx context.Context, tx *sql.Tx, eventID int64) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM event_detail WHERE id = ?`, eventID)
	if err != nil {
		log.Println("Error deleting event:", err)
		return err
	}
	return nil
}

// Insert the event slots
func (eventRepo *eventRepository) InsertEventSlots(ctx context.Context, tx *sql.Tx, eventID int64, slot model.EventSlot) error {
	_, err := tx.ExecContext(ctx, `
			INSERT INTO event_slot (event_id, start_time, end_time) 
			VALUES (?, ?, ?)`, eventID, slot.StartTime, slot.EndTime)
	if err != nil {
		log.Println("Error inserting event slot:", err)
		return err
	}
	return nil
}

// Delete the event slots
func (eventRepo *eventRepository) DeleteEventSlots(ctx context.Context, tx *sql.Tx, slotID int64) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM event_slot WHERE id = ?`, slotID)
	if err != nil {
		log.Println("Error deleting event slots:", err)
		return err

	}
	return nil
}

// Get the event slots
func (eventRepo *eventRepository) GetEventSlots(ctx context.Context, eventID int64) ([]model.EventSlot, error) {
	rows, err := eventRepo.dbConn.QueryContext(ctx, `SELECT id, start_time, end_time FROM event_slot WHERE event_id = ?`, eventID)
	if err != nil {
		log.Println("Error getting event slots:", err)
		return nil, err
	}
	defer rows.Close()

	var slots []model.EventSlot
	for rows.Next() {
		var slot model.EventSlot
		if err := rows.Scan(&slot.ID, &slot.StartTime, &slot.EndTime); err != nil {
			log.Println("Error scanning event slot:", err)
			return nil, err
		}
		slots = append(slots, slot)
	}

	return slots, nil
}

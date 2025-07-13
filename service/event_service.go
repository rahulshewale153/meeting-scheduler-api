package service

import (
	"context"
	"log"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/rahulshewale153/meeting-scheduler-api/repository"
	"github.com/rahulshewale153/meeting-scheduler-api/utils"
)

type eventService struct {
	transactionManager repository.TransactionManagerI
	eventRepo          repository.EventRepositoryI
}

func NewEventService(transactionManager repository.TransactionManagerI, eventRepo repository.EventRepositoryI) EventServiceI {
	return &eventService{
		transactionManager: transactionManager,
		eventRepo:          eventRepo,
	}
}

// InsertEvent inserts a new event into the database.
func (s *eventService) InsertEvent(ctx context.Context, createEventReq model.EventRequest) (int64, error) {
	tx, err := s.transactionManager.BeginTransaction(ctx)
	if err != nil {
		return 0, err
	}
	// Ensure that the transaction is rolled back or committed properly
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("Error rolling back transaction:", rollbackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			log.Println("Error committing transaction:", commitErr)
			return
		}
	}()

	eventID, err := s.eventRepo.InsertEvent(ctx, tx, createEventReq.Event)
	if err != nil {
		return 0, err
	}

	//insert event slot
	for _, slot := range createEventReq.ProposedSlots {
		slot.StartTime, err = utils.ConvertTimeToUTC(ctx, slot.StartTime)
		if err != nil {
			log.Println("Error converting start time to UTC:", err)
			return 0, err
		}

		slot.EndTime, err = utils.ConvertTimeToUTC(ctx, slot.EndTime)
		if err != nil {
			log.Println("Error converting end time to UTC:", err)
			return 0, err
		}

		if err = s.eventRepo.InsertEventSlots(ctx, tx, eventID, slot); err != nil {
			log.Println("Error inserting event slots:", err)
			return 0, err
		}
	}
	return eventID, nil
}

// UpdateEvent updates an existing event in the database.
func (s *eventService) UpdateEvent(ctx context.Context, updateEventReq model.EventRequest) error {
	tx, err := s.transactionManager.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	// Ensure that the transaction is rolled back or committed properly
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("Error rolling back transaction:", rollbackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			log.Println("Error committing transaction:", commitErr)
			return
		}
	}()

	if err := s.eventRepo.UpdateEvent(ctx, tx, updateEventReq.Event); err != nil {
		log.Println("Error updating event:", err)
		return err
	}

	existingMap := make(map[string]model.EventSlot)
	incomingMap := make(map[string]model.EventSlot)

	//Get existing slots
	existingSlots, err := s.eventRepo.GetEventSlots(ctx, updateEventReq.Event.ID)
	if err != nil {
		log.Println("Error getting existing event slots:", err)
		return err
	}

	for _, e := range existingSlots {
		key := utils.SlotKey(e)
		existingMap[key] = e
	}

	// insert new slots that are not in the existing slots
	for _, slot := range updateEventReq.ProposedSlots {
		key := utils.SlotKey(slot)
		incomingMap[key] = slot
		if _, ok := existingMap[key]; !ok {
			slot.StartTime, err = utils.ConvertTimeToUTC(ctx, slot.StartTime)
			if err != nil {
				log.Println("Error converting start time to UTC:", err)
				return err
			}

			slot.EndTime, err = utils.ConvertTimeToUTC(ctx, slot.EndTime)
			if err != nil {
				log.Println("Error converting end time to UTC:", err)
				return err
			}

			if err = s.eventRepo.InsertEventSlots(ctx, tx, updateEventReq.Event.ID, slot); err != nil {
				log.Println("Error inserting event slots:", err)
				return err
			}
		}
	}

	// delete slots that are not in the incoming request
	for key, oldSlot := range existingMap {
		if _, ok := incomingMap[key]; !ok {
			if err = s.eventRepo.DeleteEventSlots(ctx, tx, oldSlot.ID); err != nil {
				log.Println("Error deleting event slots:", err)
				return err
			}
		}
	}

	return nil
}

// DeleteEvent deletes an event from the database.
func (s *eventService) DeleteEvent(ctx context.Context, eventID int64) error {
	tx, err := s.transactionManager.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	// Ensure that the transaction is rolled back or committed properly
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("Error rolling back transaction:", rollbackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			log.Println("Error committing transaction:", commitErr)
			return
		}
	}()

	// Delete all event slots associated with the event
	if err := s.eventRepo.DeleteEventSlots(ctx, tx, eventID); err != nil {
		log.Println("Error deleting event slots:", err)
		return err
	}

	// Delete the event itself
	if err := s.eventRepo.DeleteEvent(ctx, tx, eventID); err != nil {
		log.Println("Error deleting event:", err)
		return err
	}

	return nil
}

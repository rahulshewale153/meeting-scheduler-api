package service

import (
	"context"
	"log"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/rahulshewale153/meeting-scheduler-api/repository"
	"github.com/rahulshewale153/meeting-scheduler-api/utils"
)

type userAvailabilityService struct {
	transactionManager   repository.TransactionManagerI
	userAvailabilityRepo repository.UserAvailabilityRepositoryI
}

func NewUserAvailabilityService(transactionManager repository.TransactionManagerI, userAvailabilityRepo repository.UserAvailabilityRepositoryI) UserAvailabilityServiceI {
	return &userAvailabilityService{
		transactionManager:   transactionManager,
		userAvailabilityRepo: userAvailabilityRepo,
	}
}

// InsertUserAvailability inserts a new user availability record into the database.
func (s *userAvailabilityService) InsertUserAvailability(ctx context.Context, userAvailability model.UserAvailability) error {
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

	userID := userAvailability.UserID
	eventID := userAvailability.EventID
	for _, slot := range userAvailability.Availability {
		_, err = s.userAvailabilityRepo.InsertUserAvailability(ctx, tx, userID, eventID, slot.StartTime, slot.EndTime)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateUserAvailability updates the availability of a user for a specific event.
func (s *userAvailabilityService) UpdateUserAvailability(ctx context.Context, userAvailability model.UserAvailability) error {
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

	existingMap := make(map[string]model.EventSlot)
	incomingMap := make(map[string]model.EventSlot)

	//Get existing user availability
	existingUserAvailability, err := s.userAvailabilityRepo.GetUserAvailability(ctx, userAvailability.EventID, userAvailability.UserID)
	if err != nil {
		log.Println("Error retrieving user availability:", err)
		return err
	}

	for _, e := range existingUserAvailability {
		key := utils.SlotKey(e)
		existingMap[key] = e
	}

	for _, slot := range userAvailability.Availability {
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
			_, err = s.userAvailabilityRepo.InsertUserAvailability(ctx, tx, userAvailability.UserID, userAvailability.EventID, slot.StartTime, slot.EndTime)
			if err != nil {
				log.Println("Error inserting user availability:", err)
				return err
			}
		}
	}

	// Delete slots that are in existing but not in incoming
	for key, oldSlot := range existingMap {
		if _, ok := incomingMap[key]; !ok {
			err = s.userAvailabilityRepo.DeleteUserAvailability(ctx, tx, userAvailability.UserID, oldSlot.ID)
			if err != nil {
				log.Println("Error deleting user availability:", err)
				return err
			}
		}
	}

	return nil
}

// DeleteUserAvailability deletes a user availability record from the database.
func (s *userAvailabilityService) DeleteUserAvailability(ctx context.Context, userID int64, eventID int64) error {
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

	err = s.userAvailabilityRepo.DeleteUserAvailability(ctx, tx, userID, eventID)
	if err != nil {
		log.Println("Error deleting user availability:", err)
		return err
	}

	return nil
}

// GetUserAvailability retrieves the availability of a specific user for a specific event.
func (s *userAvailabilityService) GetUserAvailability(ctx context.Context, eventID int64, userID int64) ([]model.EventSlot, error) {
	slots, err := s.userAvailabilityRepo.GetUserAvailability(ctx, eventID, userID)
	if err != nil {
		log.Println("Error retrieving user availability:", err)
		return nil, err
	}
	return slots, nil
}

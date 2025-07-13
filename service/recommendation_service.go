package service

import (
	"context"
	"sort"
	"time"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
	"github.com/rahulshewale153/meeting-scheduler-api/repository"
	"github.com/rahulshewale153/meeting-scheduler-api/utils"
)

type recommendationService struct {
	eventRepo            repository.EventRepositoryI
	userAvailabilityRepo repository.UserAvailabilityRepositoryI
}

// NewRecommendationService creates a new instance of recommendationService
func NewRecommendationService(eventRepo repository.EventRepositoryI, userAvailabilityRepo repository.UserAvailabilityRepositoryI) RecommendationServiceI {
	return &recommendationService{eventRepo: eventRepo, userAvailabilityRepo: userAvailabilityRepo}
}

func (s *recommendationService) GetRecommendedSlots(ctx context.Context, eventID int64) ([]model.SlotRecommendation, error) {
	results := []model.SlotRecommendation{}
	// Get the event details
	event, err := s.eventRepo.GetEvent(ctx, eventID)
	if err != nil {
		return results, err
	}

	//Get the event slot
	eventSlots, err := s.eventRepo.GetEventSlots(ctx, eventID)
	if err != nil {
		return results, err
	}

	// Get all users' availability for the event
	userAvailability, err := s.userAvailabilityRepo.GetAllEventUsers(ctx, eventID)
	if err != nil {
		return results, err
	}
	// If no users are available, return an empty slice
	if len(userAvailability) == 0 {
		return results, nil
	}

	// Step 1: Normalize event slots into time-frame
	eventSlotMap := make(map[string]model.EventSlot)
	userSlotMap := make(map[string][]int64)

	for _, es := range eventSlots {
		timeFrames := breakIntoTimeFrames(es, time.Duration(event.DurationMinutes)*time.Minute)
		for _, frame := range timeFrames {
			key := utils.SlotKey(frame)
			eventSlotMap[key] = frame
		}
	}

	// Step 2: Check user availability per time-frame
	for userID, slots := range userAvailability {
		for _, slot := range slots {
			timeFrames := breakIntoTimeFrames(slot, time.Duration(event.DurationMinutes)*time.Minute)
			for _, frame := range timeFrames {
				key := utils.SlotKey(frame)
				if _, ok := eventSlotMap[key]; ok {
					userSlotMap[key] = append(userSlotMap[key], userID)
				}
			}
		}
	}

	// Step 3: Build result
	totalUsers := make(map[int64]bool)
	for userID := range userAvailability {
		totalUsers[userID] = true
	}

	for key, users := range userSlotMap {
		slot := eventSlotMap[key]

		available := utils.Unique(users)
		unavailable := utils.Difference(totalUsers, available)

		results = append(results, model.SlotRecommendation{
			Slot:        slot,
			Available:   available,
			Unavailable: unavailable,
		})
	}

	// Step 4: Sort by number of available users descending
	sort.Slice(results, func(i, j int) bool {
		return len(results[i].Available) > len(results[j].Available)
	})

	return results, nil
}

func breakIntoTimeFrames(slot model.EventSlot, duration time.Duration) []model.EventSlot {
	var timeFrames []model.EventSlot
	start := slot.StartTime
	for start.Add(duration).Before(slot.EndTime) || start.Add(duration).Equal(slot.EndTime) {
		timeFrames = append(timeFrames, model.EventSlot{
			StartTime: start,
			EndTime:   start.Add(duration),
		})
		start = start.Add(duration)
	}
	return timeFrames
}

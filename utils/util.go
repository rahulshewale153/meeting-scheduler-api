package utils

import (
	"time"

	"github.com/rahulshewale153/meeting-scheduler-api/model"
)

func SlotKey(slot model.EventSlot) string {
	return slot.StartTime.UTC().Format(time.RFC3339) + "_" + slot.EndTime.UTC().Format(time.RFC3339)
}

func Unique(ints []int64) []int64 {
	seen := make(map[int64]struct{})
	result := []int64{}
	for _, i := range ints {
		if _, ok := seen[i]; !ok {
			seen[i] = struct{}{}
			result = append(result, i)
		}
	}
	return result
}

func Difference(all map[int64]bool, present []int64) []int64 {
	presentSet := make(map[int64]bool)
	for _, id := range present {
		presentSet[id] = true
	}

	diff := []int64{}
	for id := range all {
		if !presentSet[id] {
			diff = append(diff, id)
		}
	}
	return diff
}

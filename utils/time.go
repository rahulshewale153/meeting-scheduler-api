package utils

import (
	"context"
	"time"
)

// convert incoming time string to utc time in string format
func ConvertTimeToUTC(ctx context.Context, timeStr time.Time) (time.Time, error) {
	// Parse the incoming time string
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return time.Time{}, err
	}

	// Convert the time to UTC
	utcTime := timeStr.In(loc)
	return utcTime, nil
}

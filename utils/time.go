package utils

import (
	"context"
	"time"
)

//convert incoming time string to utc time in string format

func ConvertToUTC(ctx context.Context, timeStr string) (string, error) {
	// Parse the incoming time string
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return "", err
	}

	t, err := time.ParseInLocation(time.RFC3339, timeStr, loc)
	if err != nil {
		return "", err
	}

	// Format the time in UTC
	return t.UTC().Format(time.DateTime), nil
}

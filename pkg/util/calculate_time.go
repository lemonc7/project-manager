package util

import (
	"math"
	"time"
)

func CalculateRemainingDays(startedAt *time.Time, duration *int32) *int32 {
	if startedAt == nil || duration == nil {
		return nil
	}

	endDate := startedAt.AddDate(0, 0, int(*duration))

	remaining := time.Until(endDate)
	days := int32(math.Ceil(remaining.Hours() / 24))
	return &days
}

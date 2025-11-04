package utils

import (
	"time"

	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
)

func GroupTimesByDate(times []time.Time) map[time.Time][]time.Time {
	grouped := make(map[time.Time][]time.Time)
	for _, t := range times {
		key := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		grouped[key] = append(grouped[key], t)
	}
	return grouped
}

func IsTimeInPreferredTimes(time time.Time, prefTimes *notification.PreferredTimes) bool {
	if prefTimes == nil {
		return false
	}
	for weekday, timeRanges := range *prefTimes {
		if time.Weekday() != weekday {
			continue
		}
		if len(timeRanges) == 0 {
			return true
		}
		for _, timeRange := range timeRanges {
			if IsTimeInTimeRange(time, timeRange) {
				return true
			}
		}
	}
	return false
}

func IsTimeInTimeRange(time time.Time, timeRange subscription.TimeRange) bool {
	timeStr := time.Format("15:04")
	return timeStr >= timeRange.TimeStart && timeStr < timeRange.TimeEnd
}

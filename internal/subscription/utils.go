package subscription

var lessonTimeRange = map[int]TimeRange{
	1: {TimeStart: "08:50", TimeEnd: "10:20"},
	2: {TimeStart: "10:35", TimeEnd: "12:05"},
	3: {TimeStart: "12:35", TimeEnd: "14:05"},
	4: {TimeStart: "14:15", TimeEnd: "15:45"},
	5: {TimeStart: "15:55", TimeEnd: "17:20"},
	6: {TimeStart: "17:30", TimeEnd: "19:00"},
	7: {TimeStart: "19:10", TimeEnd: "20:30"},
	8: {TimeStart: "20:40", TimeEnd: "22:00"},
}

func lessonsToTimeRanges(lessons ...int) []TimeRange {
	ranges := make([]TimeRange, 0, len(lessons))
	for _, lesson := range lessons {
		if timeRange, exists := lessonTimeRange[lesson]; exists {
			ranges = append(ranges, timeRange)
		}
	}
	return ranges
}

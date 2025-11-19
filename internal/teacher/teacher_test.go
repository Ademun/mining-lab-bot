package teacher

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateWeekNumber(t *testing.T) {
	type testCase struct {
		currentWeek  int
		currentTime  string
		targetTime   string
		expectedWeek int
	}

	tests := []testCase{
		{currentWeek: 1, currentTime: "2025-11-19", targetTime: "2025-11-20", expectedWeek: 1},
		{currentWeek: 1, currentTime: "2025-11-19", targetTime: "2025-11-24", expectedWeek: 2},
		{currentWeek: 1, currentTime: "2025-11-19", targetTime: "2025-11-30", expectedWeek: 2},
		{currentWeek: 1, currentTime: "2025-11-19", targetTime: "2025-12-01", expectedWeek: 1},
		{currentWeek: 1, currentTime: "2025-12-01", targetTime: "2025-12-02", expectedWeek: 1},
		{currentWeek: 1, currentTime: "2025-12-01", targetTime: "2025-12-08", expectedWeek: 2},
		{currentWeek: 1, currentTime: "2025-11-17", targetTime: "2025-11-17", expectedWeek: 1},
		{currentWeek: 1, currentTime: "2025-11-17", targetTime: "2025-11-23", expectedWeek: 1},
		{currentWeek: 1, currentTime: "2025-11-17", targetTime: "2025-11-24", expectedWeek: 2},
		{currentWeek: 2, currentTime: "2025-11-19", targetTime: "2025-11-20", expectedWeek: 2},
		{currentWeek: 2, currentTime: "2025-11-19", targetTime: "2025-11-24", expectedWeek: 1},
		{currentWeek: 2, currentTime: "2025-11-19", targetTime: "2025-12-01", expectedWeek: 2},
		{currentWeek: 2, currentTime: "2025-12-01", targetTime: "2025-12-02", expectedWeek: 2},
		{currentWeek: 2, currentTime: "2025-12-01", targetTime: "2025-12-08", expectedWeek: 1},
		{currentWeek: 1, currentTime: "2025-11-19", targetTime: "2025-12-15", expectedWeek: 1},
		{currentWeek: 1, currentTime: "2025-11-19", targetTime: "2026-01-05", expectedWeek: 2},
		{currentWeek: 1, currentTime: "2025-11-19", targetTime: "2026-01-05", expectedWeek: 2},
		{currentWeek: 2, currentTime: "2025-11-19", targetTime: "2025-12-15", expectedWeek: 2},
		{currentWeek: 2, currentTime: "2025-11-19", targetTime: "2026-01-05", expectedWeek: 1},
		{currentWeek: 1, currentTime: "2025-11-19", targetTime: "2026-01-19", expectedWeek: 2},
		{currentWeek: 1, currentTime: "2025-12-15", targetTime: "2025-12-20", expectedWeek: 1},
		{currentWeek: 1, currentTime: "2025-12-15", targetTime: "2025-12-22", expectedWeek: 2},
		{currentWeek: 2, currentTime: "2026-01-05", targetTime: "2026-01-12", expectedWeek: 1},
	}

	for i, tCase := range tests {
		t.Run(fmt.Sprintf("test_calc_week_number_%d", i), func(t *testing.T) {
			currTime, _ := time.Parse("2006-01-02", tCase.currentTime)
			targetTime, _ := time.Parse("2006-01-02", tCase.targetTime)
			actual := calculateWeekNumber(tCase.currentWeek, currTime, targetTime)
			assert.Equal(t, tCase.expectedWeek, actual)
		})
	}
}

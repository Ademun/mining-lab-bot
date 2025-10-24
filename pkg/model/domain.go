package model

import (
	"fmt"
	"strings"
	"time"
)

type LabType int

const (
	LabPerformance LabType = iota
	LabDefence
)

func (t LabType) String() string {
	switch t {
	case LabPerformance:
		return "Выполнение"
	case LabDefence:
		return "Защита"
	}
	return "Unknown"
}

type Slot struct {
	ID            int
	LabName       string
	LabNumber     int
	LabAuditorium int
	LabOrder      int
	LabType       LabType
	Available     []TimeTeachers
	URL           string
}

type TimeTeachers struct {
	Time     time.Time
	Teachers []Teacher
}

func (s Slot) Key() string {
	var timeKey strings.Builder
	for _, available := range s.Available {
		timeKey.WriteString(available.Time.Format(time.RFC3339))
	}
	return fmt.Sprintf("%d_%s", s.ID, timeKey.String())
}

type Subscription struct {
	UUID          string        `db:"uuid"`
	UserID        int           `db:"user_id"`
	LabNumber     int           `db:"lab_number"`
	LabAuditorium int           `db:"lab_auditorium"`
	Weekday       *time.Weekday `db:"weekday"`
	DayTime       *string       `db:"day_time"`
}

type Notification struct {
	UserID         int
	PreferredTimes []PreferredTime
	Slot           Slot
}

type PreferredTime struct {
	Weekday time.Weekday
	DayTime string
}

type Teacher struct {
	ID         int          `db:"id"`
	Name       string       `db:"name"`
	Auditorium int          `db:"auditorium"`
	WeekNumber int          `db:"week_number"`
	Weekday    time.Weekday `db:"weekday"`
	TimeStart  string       `db:"time_start"`
	TimeEnd    string       `db:"time_end"`
}

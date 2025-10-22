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
	UUID          string
	UserID        int
	ChatID        int
	LabNumber     int
	LabAuditorium int
	Weekday       *time.Weekday
	DayTime       string
}

type Notification struct {
	UserID        int
	ChatID        int
	PreferredTime time.Time
	Slot          Slot
}

type Teacher struct {
	ID         int
	Name       string
	Auditorium int
	WeekNumber int
	Weekday    time.Weekday
	TimeStart  string
	TimeEnd    string
}

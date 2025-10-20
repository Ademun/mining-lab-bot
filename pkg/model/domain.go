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
	LabType       LabType
	Available     []time.Time
	URL           string
}

func (s Slot) Key() string {
	var timeKey strings.Builder
	for _, dateTime := range s.Available {
		timeKey.WriteString(dateTime.Format(time.RFC3339))
	}
	return fmt.Sprintf("%d_%s", s.ID, timeKey.String())
}

type Subscription struct {
	UUID          string
	UserID        int
	ChatID        int
	LabNumber     int
	LabAuditorium int
	Weekday       time.Weekday
	DayTime       string
	Teacher       string
}

type Notification struct {
	UserID        int
	ChatID        int
	PreferredTime time.Time
	Slot          Slot
}

type Teacher struct {
	UUID       string
	Name       string
	Weekday    time.Weekday
	TimeStart  string
	TimeEnd    string
	Auditorium int
	WeekNumber int
	Difficulty int
}

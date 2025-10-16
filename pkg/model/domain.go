package model

import (
	"fmt"
	"strings"
	"time"
)

type Slot struct {
	ID            int
	LabNumber     int
	LabName       string
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

type Subscription struct {
	UUID          string
	UserID        int
	ChatID        int
	LabNumber     int
	LabAuditorium int
}

type Notification struct {
	UserID int
	ChatID int
	Slot   Slot
}

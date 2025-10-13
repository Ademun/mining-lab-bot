package model

import (
	"time"
)

type Slot struct {
	ID            int
	LabNumber     int
	LabName       string
	LabAuditorium int
	LabType       LabType
	DateTime      time.Time
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

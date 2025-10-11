package cmd

import (
	"fmt"
	"time"
)

func formatDateTime(t time.Time) string {
	now := time.Now()
	date := formatDateRelative(t, now)
	timeStr := t.Format("15:04")

	return fmt.Sprintf("%s в %s", date, timeStr)
}

func formatDateRelative(t time.Time, now time.Time) string {
	targetDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	daysDiff := int(targetDay.Sub(todayStart).Hours() / 24)

	switch daysDiff {
	case 0:
		return "Сегодня"
	case 1:
		return "Завтра"
	case 2:
		return "Послезавтра"
	default:
		return formatDateLong(t)
	}
}

func formatDateLong(t time.Time) string {
	months := []string{
		"января", "февраля", "марта", "апреля", "мая", "июня",
		"июля", "августа", "сентября", "октября", "ноября", "декабря",
	}

	return fmt.Sprintf("%d %s", t.Day(), months[t.Month()-1])
}

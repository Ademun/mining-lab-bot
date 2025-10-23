package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

func formatAvailableSlots(slots []model.TimeTeachers) map[string][]model.TimeTeachers {
	slotsMap := make(map[string][]model.TimeTeachers)
	for _, slot := range slots {
		key := slot.Time.Format("2006-01-02")
		slotsMap[key] = append(slotsMap[key], slot)
	}
	return slotsMap
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

func formatDuration(d time.Duration) string {
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second

	if days > 0 {
		return fmt.Sprintf("%dд %dч %dм", days, hours, minutes)
	}
	return fmt.Sprintf("%dч %dм %dс", hours, minutes, seconds)
}

func formatPollingMode(mode config.PollingMode) string {
	switch mode {
	case config.ModeNormal:
		return "стандартный"
	case config.ModeAggressive:
		return "агрессивный"
	default:
		return "неизвестный"
	}
}

func parseTime(input string) (time.Time, error) {
	re := regexp.MustCompile(`^(\d{1,2})[:.]?(\d{2})$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return time.Now(), fmt.Errorf("неверный формат времени")
	}
	hours, err := strconv.Atoi(matches[1])
	if err != nil || hours < 0 || hours > 23 {
		return time.Now(), fmt.Errorf("часы должны быть в диапазоне 0-23")
	}
	minutes, err := strconv.Atoi(matches[2])
	if err != nil || minutes < 0 || minutes > 59 {
		return time.Now(), fmt.Errorf("минуты должны быть в диапазоне 0-59")
	}

	now := time.Now()
	parsedTime := time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, 0, 0,
		now.Location())

	return parsedTime, nil
}

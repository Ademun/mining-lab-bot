package utils

import (
	"fmt"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/config"
)

var WeekdayLocale = map[int]string{
	0: "Воскресенье",
	1: "Понедельник",
	2: "Вторник",
	3: "Среда",
	4: "Четверг",
	5: "Пятница",
	6: "Суббота",
}

var TimeStartToLongLessonTime = map[string]string{
	"08:50": "08:50 - 10:20 - 1️⃣ пара",
	"10:35": "10:35 - 12:05 - 2️⃣ пара",
	"12:35": "12:35 - 14:05 - 3️⃣ пара",
	"14:15": "14:15 - 15:45 - 4️⃣ пара",
	"15:55": "15:55 - 17:20 - 5️⃣ пара",
	"17:30": "17:30 - 19:00 - 6️⃣ пара",
	"19:10": "19:10 - 20:30 - 7️⃣ пара",
	"20:40": "20:40 - 22:00 - 8️⃣ пара",
}

var TimeStartToShortLessonTime = map[string]string{
	"08:50": "1️⃣ 08:50 - 10:20",
	"10:35": "2️⃣ 10:35 - 12:05",
	"12:35": "3️⃣ 12:35 - 14:05",
	"14:15": "4️⃣ 14:15 - 15:45",
	"15:55": "5️⃣ 15:55 - 17:20",
	"17:30": "6️⃣ 17:30 - 19:00",
	"19:10": "7️⃣ 19:10 - 20:30",
	"20:40": "8️⃣ 20:40 - 22:00",
}

var Months = []string{
	"января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

func FormatDateRelative(t time.Time, now time.Time) string {
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
		return FormatDateLong(t)
	}
}

func FormatDateLong(t time.Time) string {
	return fmt.Sprintf("%d %s (%s)", t.Day(), Months[t.Month()-1], WeekdayLocale[int(t.Weekday())])
}

func FormatDuration(d time.Duration) string {
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

func FormatPollingMode(mode config.PollingMode) string {
	switch mode {
	case config.ModeNormal:
		return "стандартный"
	case config.ModeAggressive:
		return "агрессивный"
	default:
		return "неизвестный"
	}
}

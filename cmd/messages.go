package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/metrics"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

func startMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>👋 Привет!</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>Я бот для записи на лабораторные работы</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>Буду следить за появлением доступных записей и сразу уведомлять тебя, когда появится нужная </b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>Используй /help для просмотра доступных команд</b>")
	return sb.String()
}

func helpMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>📖 Справка</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>📝 Подписка:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>/sub - создать подписку</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>⚙️ Управление:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>/unsub - удалить подписку</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>/list - посмотреть подписки</b>")
	return sb.String()
}

func subAskLabNumberMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>📚 Введите номер лабораторной работы</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Например: 7")
	return sb.String()
}

func subAskAuditoriumMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>🚪 Введите номер аудитории</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Например: 233")
	return sb.String()
}

func subAskWeekdayMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>📅 Выберите день недели</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Или пропустите, если день не важен")
	return sb.String()
}

func subAskLessonMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>🕐 Выбери время</b>")
	return sb.String()
}

var weekDayLocale = map[int]string{
	0: "Воскресенье",
	1: "Понедельник",
	2: "Вторник",
	3: "Среда",
	4: "Четверг",
	5: "Пятница",
	6: "Суббота",
}

func subConfirmationMessage(data *subscriptionData) string {
	labNumber := data.LabNumber
	auditorium := data.LabAuditorium
	weekday := data.Weekday
	timeStr := data.Daytime

	var sb strings.Builder
	sb.WriteString("<b>✅ Создать подписку?</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>📚 Лаба:</b> %d", labNumber))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>🚪 Аудитория:</b> %d", auditorium))

	if weekday != nil {
		sb.WriteString(repeatLineBreaks(2))
		sb.WriteString(fmt.Sprintf("<b>📅 День:</b> %s", weekDayLocale[int(*weekday)]))
	}

	if timeStr != "" {
		sb.WriteString(repeatLineBreaks(2))
		sb.WriteString(fmt.Sprintf("<b>🕐 Время:</b> %s", timeLessonMap[timeStr]))
	}

	return sb.String()
}

func subLabNumberValidationErrorMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Номер лабы должен быть числом в диапазоне от 1 до 999</b>")
	return sb.String()
}

func subAuditoriumNumberValidationErrorMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Номер аудитории должен быть числом в диапазоне от 1 до 999</b>")
	return sb.String()
}

func subCancelledMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Создание подписки отменено</b>")
	return sb.String()
}

func subCreationErrorMessage(err error) string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Произошла ошибка при создании подписки:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>%s</b>", err.Error()))
	return sb.String()
}

func subCreationSuccessMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>✅ Подписка создана!</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🔔 Вы получите уведомление, когда появится нужная запись</b>")
	return sb.String()
}

func unsubEmptyListMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>🔍 У вас нет подписок на лабы</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Используйте команду /sub для создания подписки")
	return sb.String()
}

func unsubSelectMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>🗑️ Выберите подписку для удаления:</b>")
	return sb.String()
}

func unsubConfirmDeleteAllMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>⚠️ Удалить все подписки?</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Это действие нельзя отменить")
	return sb.String()
}

func unsubDeleteAllSuccessMessage(count int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>✅ Удалено подписок: %d</b>", count))
	return sb.String()
}

func subsFetchingErrorMessage(err error) string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Произошла ошибка при получении списка подписок:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>%s</b>", err.Error()))
	return sb.String()
}

func unsubErrorMessage(err error) string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Произошла ошибка при отписке:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>%s</b>", err.Error()))
	return sb.String()
}

func unsubSuccessMessage(labNumber, labAuditorium int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>✅ Вы больше не подписаны на лабу №%d в ауд. №%d</b>",
		labNumber, labAuditorium))
	return sb.String()
}

func listEmptySubsMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>🔍 У вас нет подписок на лабы</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Используйте команду /sub для создания подписки")
	return sb.String()
}

func listSubsSuccessMessage(subs []model.Subscription) string {
	var sb strings.Builder
	sb.WriteString("<b>📋 Ваши подписки:</b>")
	sb.WriteString(repeatLineBreaks(2))
	for idx, sub := range subs {
		sb.WriteString(fmt.Sprintf("<b>%d.</b> Лаба №%d, ауд. №%d", idx+1,
			sub.LabNumber, sub.LabAuditorium))
		if idx == len(subs)-1 {
			break
		}
		sb.WriteString(repeatLineBreaks(2))
	}
	return sb.String()
}

func permissionDeniedErrorMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Доступ запрещён. Команда доступна только для разработчика</b>")
	return sb.String()
}

func statsSuccessMessage(snapshot *metrics.Metrics) string {
	uptime := time.Since(snapshot.StartTime)
	var sb strings.Builder
	sb.WriteString("<b>📊 Статистика сервиса</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🕐 Общее время работы:</b> ")
	sb.WriteString(formatDuration(uptime))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🔍 Опросы:</b>")
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Всего опросов: <b>%d</b>",
		snapshot.PollingMetrics.TotalPolls))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Режим: <b>%s</b>",
		formatPollingMode(snapshot.PollingMetrics.Mode)))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Ошибки парсинга: <b>%d</b>",
		snapshot.PollingMetrics.ParsingErrors))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Ошибки получения: <b>%d</b>",
		snapshot.PollingMetrics.FetchErrors))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Среднее время опроса: <b>%s</b>",
		snapshot.PollingMetrics.AveragePollingTime.Round(time.Millisecond)))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Среднее количество слотов: <b>%d</b>",
		snapshot.PollingMetrics.AverageSlotNumber))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🔔 Уведомления:</b>")
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Всего уведомлений: <b>%d</b>",
		snapshot.NotificationMetrics.TotalNotifications))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Размер кеша: <b>%d</b>",
		snapshot.NotificationMetrics.CacheLength))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Среднее количество уведомлений: <b>%d</b>",
		snapshot.NotificationMetrics.AverageNotifications))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>📝 Подписки:</b>")
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Активных подписок: <b>%d</b>",
		snapshot.SubscriptionMetrics.TotalSubscriptions))
	return sb.String()
}

var timeLessonMap = map[string]string{
	"08:50": "1️⃣ 08:50 - 10:20",
	"10:35": "2️⃣ 10:35 - 12:05",
	"12:35": "3️⃣ 12:35 - 14:05",
	"14:15": "4️⃣ 14:15 - 15:45",
	"15:55": "5️⃣ 15:55 - 17:20",
	"17:30": "6️⃣ 17:30 - 19:00",
	"19:10": "7️⃣ 19:10 - 20:30",
	"20:40": "8️⃣ 20:40 - 22:00",
}

func notifySuccessMessage(notif *model.Notification) string {
	slot := &notif.Slot
	var sb strings.Builder
	sb.WriteString("<b>🔥 Появилась запись!</b>")
	sb.WriteString(repeatLineBreaks(3))
	longName := slot.LabName
	if slot.LabOrder != 0 {
		longName += fmt.Sprintf(" (%d-ое место)", slot.LabOrder)
	}
	sb.WriteString(fmt.Sprintf("<b>📚 Лаба №%d. %s</b>", slot.LabNumber, longName))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>🚪 Аудитория №%d</b>", slot.LabAuditorium))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🗓️ Когда:</b>")
	sb.WriteString(repeatLineBreaks(1))
	writeSlotsInfo(slot, &sb, notif.PreferredTime)
	return sb.String()
}

func writeSlotsInfo(slot *model.Slot, sb *strings.Builder, preferredTime model.PreferredTime) {
	available := formatAvailableSlots(slot.Available)

	keys := sortDatesByPreference(available, preferredTime)

	for idx, k := range keys {
		val := available[k]
		parsedTime, _ := time.Parse("2006-01-02", k)
		relativeDate := formatDateRelative(parsedTime, time.Now())

		isPreferredDate := parsedTime.Weekday() == preferredTime.Weekday

		sb.WriteString(fmt.Sprintf("<b>⠀⠀%s:</b>", relativeDate))
		sb.WriteString(repeatLineBreaks(1))

		sortedSlots := sortSlotsByPreference(val, preferredTime.DayTime, isPreferredDate)

		for idx, v := range sortedSlots {
			timeStart := v.Time.Format("15:04")
			timePart := timeLessonMap[timeStart]
			teacherPart := make([]string, len(v.Teachers))
			for idx, teacher := range v.Teachers {
				teacherPart[idx] = teacher.Name
			}

			isPreferredSlot := isPreferredDate && timeStart == preferredTime.DayTime
			if isPreferredSlot {
				sb.WriteString(fmt.Sprintf("<b>⠀⠀%s %s ⭐ Ваше время</b>", timePart, strings.Join(teacherPart, ", ")))
			} else {
				sb.WriteString(fmt.Sprintf("<b>⠀⠀%s %s</b>", timePart, strings.Join(teacherPart, ", ")))
			}

			if idx != len(sortedSlots)-1 {
				sb.WriteString(repeatLineBreaks(1))
			}
		}
		if idx != len(keys)-1 {
			sb.WriteString(repeatLineBreaks(2))
		}
	}
}

func sortDatesByPreference(available map[string][]model.TimeTeachers, preferredTime model.PreferredTime) []string {
	keys := make([]string, 0, len(available))
	for k := range available {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", keys[i])
		dateJ, _ := time.Parse("2006-01-02", keys[j])

		isPreferredI := dateI.Weekday() == preferredTime.Weekday
		isPreferredJ := dateJ.Weekday() == preferredTime.Weekday

		if isPreferredI && !isPreferredJ {
			return true
		}
		if !isPreferredI && isPreferredJ {
			return false
		}

		return dateI.Before(dateJ)
	})

	return keys
}

func sortSlotsByPreference(slots []model.TimeTeachers, preferredDayTime string, isPreferredDate bool) []model.TimeTeachers {
	sorted := make([]model.TimeTeachers, len(slots))
	copy(sorted, slots)

	if !isPreferredDate {
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Time.Before(sorted[j].Time)
		})
		return sorted
	}

	sort.Slice(sorted, func(i, j int) bool {
		timeI := sorted[i].Time.Format("15:04")

		timeJ := sorted[j].Time.Format("15:04")

		isPreferredI := timeI == preferredDayTime
		isPreferredJ := timeJ == preferredDayTime

		// Предпочтительное время всегда первое
		if isPreferredI && !isPreferredJ {
			return true
		}
		if !isPreferredI && isPreferredJ {
			return false
		}

		// Остальные по времени
		return sorted[i].Time.Before(sorted[j].Time)
	})

	return sorted
}

func repeatLineBreaks(breaks int) string {
	var sb strings.Builder
	for range breaks {
		sb.WriteString("\n")
	}
	return sb.String()
}

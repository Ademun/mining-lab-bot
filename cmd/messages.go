package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
)

func helpMsg() string {
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

func startMsg() string {
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

func askLabTypeMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>📝 Выберите тип лабораторной работы")
	return sb.String()
}

func askLabNumberMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>📚 Введите номер лабораторной работы</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Например: 7")
	return sb.String()
}

func labNumberValidationErrorMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Номер лабы должен быть числом в диапазоне от 1 до 999</b>")
	return sb.String()
}

func askLabAuditoriumMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>🚪 Введите номер аудитории</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Например: 233")
	return sb.String()
}

func labAuditoriumValidationErrorMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Номер аудитории должен быть числом в диапазоне от 1 до 999</b>")
	return sb.String()
}

func askLabDomainMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>⚛️ Выберите вид лабораторной работы")
	return sb.String()
}

func askLabWeekdayMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>📅 Выберите день недели</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Или пропустите, если день не важен")
	return sb.String()
}

func askLabLessonsMsg() string {
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

func askSubCreationConfirmationMsg(sub *subscription.RequestSubscription) string {
	var sb strings.Builder
	sb.WriteString("<b>✅ Создать подписку?</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>📚 Лаба: %d. %s</b>", sub.LabNumber, sub.Type.String()))
	sb.WriteString(repeatLineBreaks(2))
	if sub.LabAuditorium != nil {
		sb.WriteString(fmt.Sprintf("<b>🚪 Аудитория:</b> %d", sub.LabAuditorium))
	} else if sub.LabDomain != nil {
		sb.WriteString(fmt.Sprintf("<b>⚛️ %s</b>", sub.LabDomain))
	}
	sb.WriteString(repeatLineBreaks(2))

	if sub.Weekday != nil {
		sb.WriteString(repeatLineBreaks(2))
		sb.WriteString(fmt.Sprintf("<b>📅 День:</b> %s", weekDayLocale[*sub.Weekday]))
	}

	if sub.Lessons != nil {
		sb.WriteString(repeatLineBreaks(2))
		sb.WriteString(fmt.Sprintf("<b>🕐 Время:</b>"))
		for _, lesson := range sub.Lessons {
			sb.WriteString(fmt.Sprintf("<b>%s</b>", defaultLessons[lesson-1].Text))
		}
	}

	return sb.String()
}

func subCreationCancelledMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Создание подписки отменено</b>")
	return sb.String()
}

func subCreationErrorMsg(err error) string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Произошла ошибка при создании подписки:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>%s</b>", err.Error()))
	return sb.String()
}

func subCreationSuccessMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>✅ Подписка создана!</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🔔 Вы получите уведомление, когда появится нужная запись</b>")
	return sb.String()
}

var timeStartToLesson = map[string]string{
	"08:50": "08:50 - 10:20 - 1️⃣ пара",
	"10:35": "10:35 - 12:05 - 2️⃣ пара",
	"12:35": "12:35 - 14:05 - 3️⃣ пара",
	"14:15": "14:15 - 15:45 - 4️⃣ пара",
	"15:55": "15:55 - 17:20 - 5️⃣ пара",
	"17:30": "17:30 - 19:00 - 6️⃣ пара",
	"19:10": "19:10 - 20:30 - 7️⃣ пара",
	"20:40": "20:40 - 22:00 - 8️⃣ пара",
}

func viewSubResponseMsg(sub *subscription.ResponseSubscription) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>📚 Лаба: %d. %s</b>", sub.LabNumber, sub.LabType.String()))
	sb.WriteString(repeatLineBreaks(2))
	if sub.LabAuditorium != nil {
		sb.WriteString(fmt.Sprintf("<b>🚪 Аудитория:</b> %d", sub.LabAuditorium))
	} else if sub.LabDomain != nil {
		sb.WriteString(fmt.Sprintf("<b>⚛️ %s</b>", sub.LabDomain))
	}
	sb.WriteString(repeatLineBreaks(2))

	if sub.Weekday != nil {
		sb.WriteString(repeatLineBreaks(2))
		sb.WriteString(fmt.Sprintf("<b>📅 День:</b> %s", weekDayLocale[*sub.Weekday]))
	}

	if sub.PreferredTimes != nil {
		sb.WriteString(repeatLineBreaks(2))
		sb.WriteString(fmt.Sprintf("<b>🕐 Время:</b>"))
		for _, prefTime := range sub.PreferredTimes {
			sb.WriteString(fmt.Sprintf("<b>%s</b>", timeStartToLesson[prefTime.TimeStart]))
		}
	}
	return sb.String()
}

func emptySubsListMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>🔍 У вас нет подписок на лабы</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Используйте команду /sub для создания подписки")
	return sb.String()
}

func permissionDeniedErrorMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Доступ запрещён. Команда доступна только для разработчика</b>")
	return sb.String()
}

func statsMsg(snapshot *metrics.Metrics) string {
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
	sb.WriteString(fmt.Sprintf("  Последнее время опроса: <b>%s</b>",
		snapshot.PollingMetrics.LastPollingTime.Round(time.Millisecond)))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Количество слотов: <b>%d</b>",
		snapshot.PollingMetrics.LastSlotNumber))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("	 Количество айдишников сервиса <b>%d</b>", snapshot.PollingMetrics.LastIDNumber))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🔔 Уведомления:</b>")
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Всего уведомлений: <b>%d</b>",
		snapshot.NotificationMetrics.TotalNotifications))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Размер кеша: <b>%d</b>",
		snapshot.NotificationMetrics.CacheLength))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>📝 Подписки:</b>")
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Активных подписок: <b>%d</b>",
		snapshot.SubscriptionMetrics.TotalSubscriptions))
	return sb.String()
}

func genericServiceErrorMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Произошла ошибка сервиса. Попробуйте позже</b>")
	return sb.String()
}

func unsubSuccessMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>✅ Вы больше не подписаны на эту лабу</b>")
	return sb.String()
}

func notifySuccessMessage(notif *notification.Notification) string {
	slot := &notif.Slot
	var sb strings.Builder
	sb.WriteString("<b>🔥 Появилась запись!</b>")
	sb.WriteString(repeatLineBreaks(3))
	longName := slot.Name
	if slot.Order != 0 {
		longName += fmt.Sprintf(" (%d-ое место)", slot.Order)
	}
	sb.WriteString(fmt.Sprintf("<b>📚 Лаба №%d. %s</b>", slot.Number, longName))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>🚪 Аудитория №%d</b>", slot.Auditorium))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🗓️ Когда:</b>")
	sb.WriteString(repeatLineBreaks(1))
	writeSlotsInfo(slot, &sb, notif.PreferredTimes)
	return sb.String()
}

func writeSlotsInfo(slot *polling.Slot, sb *strings.Builder, preferredTimes notification.PreferredTimes) {
	// Группируем слоты по датам
	available := groupSlotsByDate(slot.TimesTeachers)

	keys := sortDatesByPreference(available, preferredTimes)

	// Создаём set предпочитаемых слотов (weekday + время)
	preferredSet := buildPreferredSet(preferredTimes)

	for _, k := range keys {
		teachers := available[k]
		parsedDate, _ := time.Parse("2006-01-02", k)
		relativeDate := formatDateRelative(parsedDate, time.Now())

		sb.WriteString(fmt.Sprintf("<b>⠀⠀%s:</b>", relativeDate))
		sb.WriteString(repeatLineBreaks(1))

		sortedSlots := sortSlotsByPreference(teachers, preferredTimes, parsedDate.Weekday())

		for idx, slotInfo := range sortedSlots {
			timeStart := slotInfo.Time.Format("15:04")
			timePart := timeStartToLesson[timeStart]

			preferredKey := fmt.Sprintf("%d_%s", parsedDate.Weekday(), timeStart)
			isPreferredSlot := preferredSet[preferredKey]

			if isPreferredSlot {
				sb.WriteString(fmt.Sprintf("<b>⠀⠀%s %s ⭐️Ваше время</b>", timePart, strings.Join(slotInfo.Teachers, ", ")))
			} else {
				sb.WriteString(fmt.Sprintf("<b>⠀⠀%s %s</b>", timePart, strings.Join(slotInfo.Teachers, ", ")))
			}

			if idx != len(sortedSlots)-1 {
				sb.WriteString(repeatLineBreaks(1))
			}
		}
		sb.WriteString(repeatLineBreaks(2))
	}
}

// SlotInfo содержит информацию о конкретном временном слоте
type SlotInfo struct {
	Time     time.Time
	Teachers []string
}

// groupSlotsByDate группирует слоты по датам в формате "2006-01-02"
func groupSlotsByDate(timesTeachers map[time.Time][]string) map[string][]SlotInfo {
	grouped := make(map[string][]SlotInfo)

	for t, teachers := range timesTeachers {
		dateKey := t.Format("2006-01-02")
		grouped[dateKey] = append(grouped[dateKey], SlotInfo{
			Time:     t,
			Teachers: teachers,
		})
	}

	return grouped
}

// buildPreferredSet создаёт set предпочитаемых временных слотов
func buildPreferredSet(preferredTimes notification.PreferredTimes) map[string]bool {
	preferredSet := make(map[string]bool)

	for weekday, timeRanges := range preferredTimes {
		for _, tr := range timeRanges {
			// Проверяем, попадает ли конкретное время в диапазон
			// Сохраняем начало диапазона как ключ
			key := fmt.Sprintf("%d_%s", weekday, tr.TimeStart)
			preferredSet[key] = true
		}
	}

	return preferredSet
}

// isTimeInRange проверяет, попадает ли время в один из предпочитаемых диапазонов
func isTimeInRange(timeStr string, timeRanges []subscription.TimeRange) bool {
	for _, tr := range timeRanges {
		if timeStr >= tr.TimeStart && timeStr <= tr.TimeEnd {
			return true
		}
	}
	return false
}

func sortDatesByPreference(available map[string][]SlotInfo, preferredTimes notification.PreferredTimes) []string {
	keys := make([]string, 0, len(available))
	for k := range available {
		keys = append(keys, k)
	}

	preferredWeekdays := make(map[time.Weekday]bool)
	for weekday := range preferredTimes {
		preferredWeekdays[weekday] = true
	}

	sort.Slice(keys, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", keys[i])
		dateJ, errJ := time.Parse("2006-01-02", keys[j])

		// Если ошибка парсинга, отправляем в конец
		if errI != nil {
			return false
		}
		if errJ != nil {
			return true
		}

		isPreferredI := preferredWeekdays[dateI.Weekday()]
		isPreferredJ := preferredWeekdays[dateJ.Weekday()]

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

func sortSlotsByPreference(slots []SlotInfo, preferredTimes notification.PreferredTimes, dateWeekday time.Weekday) []SlotInfo {
	sorted := make([]SlotInfo, len(slots))
	copy(sorted, slots)

	// Получаем предпочитаемые диапазоны времени для этого дня недели
	timeRanges, hasPreferences := preferredTimes[dateWeekday]

	if !hasPreferences || len(timeRanges) == 0 {
		// Если нет предпочтений, просто сортируем по времени
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Time.Before(sorted[j].Time)
		})
		return sorted
	}
	sort.Slice(sorted, func(i, j int) bool {
		timeI := sorted[i].Time.Format("15:04")
		timeJ := sorted[j].Time.Format("15:04")

		isPreferredI := isTimeInRange(timeI, timeRanges)
		isPreferredJ := isTimeInRange(timeJ, timeRanges)

		// Предпочитаемые слоты идут первыми
		if isPreferredI && !isPreferredJ {
			return true
		}
		if !isPreferredI && isPreferredJ {
			return false
		}

		// Если оба предпочитаемые или оба нет - сортируем по времени
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

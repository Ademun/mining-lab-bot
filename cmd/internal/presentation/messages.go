package presentation

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Ademun/mining-lab-bot/cmd/internal/utils"
	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
)

func HelpCmdMsg() string {
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

func StartCmdMsg() string {
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

// Feedback text flow

func FeedbackCmdMsg() string {
	return "<b>🖊️ Напишите ваши пожелания и идеи</b>"
}

func FeedbackRedirectMsg(userID int64, feedback string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>От пользователя: %d</b>", userID))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(feedback)
	return sb.String()
}

func FeedbackReplyMsg() string {
	return "<b>😊 Спасибо за ваше предложение! Оно будет принято к рассмотрению</b>"
}

// ===

func GenericServiceErrorMsg() string {
	return "<b>❌ Произошла ошибка сервиса. Попробуйте позже</b>"
}

func ValidationErrorMsg(cause string) string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Ошибка валидации:</b>")
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("<b>%s</b>", cause))
	return sb.String()
}

// Subscription creation flow

func AskLabTypeMsg() string {
	return "<b>📝 Выберите тип лабораторной работы</b>"
}

func AskLabNumberMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>📚 Введите номер лабораторной работы</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Например: 7")
	return sb.String()
}

func AskLabAuditoriumMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>🚪 Введите номер аудитории</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Например: 233")
	return sb.String()
}

func AskLabDomainMsg() string {
	return "<b>⚛️ Выберите вид лабораторной работы</b>"
}

func AskWeekdayMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>📅 Выберите день недели</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Или пропустите, если день не важен")
	return sb.String()
}

func AskLessonsMsg(lessons []int) string {
	var sb strings.Builder
	sb.WriteString("<b>🕐 Выбери время</b>")
	if len(lessons) > 0 {
		sb.WriteString(repeatLineBreaks(2))
		sb.WriteString("<b>Выбранные пары:</b>")
		slices.Sort(lessons)
		for _, lesson := range lessons {
			sb.WriteString(repeatLineBreaks(1))
			sb.WriteString(fmt.Sprintf("<b>%s</b>", utils.LessonNumberToLessonName[lesson]))
		}
	}
	return sb.String()
}

func AskSubCreationConfirmationMsg(sub *subscription.RequestSubscription) string {
	var sb strings.Builder
	sb.WriteString("<b>✅ Создать подписку?</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>📚 Лаба: %d. %s</b>", sub.LabNumber, sub.Type.String()))
	sb.WriteString(repeatLineBreaks(2))
	if sub.LabAuditorium != nil {
		sb.WriteString(fmt.Sprintf("<b>🚪 Аудитория:</b> %d", *sub.LabAuditorium))
	} else if sub.LabDomain != nil {
		sb.WriteString(fmt.Sprintf("<b>⚛️ %s</b>", sub.LabDomain))
	}
	sb.WriteString(repeatLineBreaks(2))

	if sub.Weekday != nil {
		sb.WriteString(fmt.Sprintf("<b>📅 День:</b> %s", utils.WeekdayLocale[*sub.Weekday]))
		sb.WriteString(repeatLineBreaks(2))
	}

	if sub.Lessons != nil {
		sb.WriteString(fmt.Sprintf("<b>🕐 Время:</b>"))
		sb.WriteString(repeatLineBreaks(2))
		for _, lesson := range sub.Lessons {
			sb.WriteString(fmt.Sprintf("<b>%s</b>", utils.DefaultLessons[lesson-1].Text))
			sb.WriteString(repeatLineBreaks(1))
		}
	}

	return sb.String()
}

func SubCreationCancelledMsg() string {
	return "<b>❌ Создание подписки отменено</b>"
}

func SubCreationSuccessMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>✅ Подписка создана!</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🔔 Вы получите уведомление, когда появится нужная запись</b>")
	return sb.String()
}

// ===

// Subscription listing flow

func EmptySubListMsg() string {
	var sb strings.Builder
	sb.WriteString("<b>🔍 У вас нет подписок на лабы</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Используйте команду /sub для создания подписки")
	return sb.String()
}

func SubViewMsg(sub *subscription.ResponseSubscription) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>📚 Лаба: %d. %s</b>", sub.LabNumber, sub.LabType.String()))
	sb.WriteString(repeatLineBreaks(2))
	if sub.LabAuditorium != nil {
		sb.WriteString(fmt.Sprintf("<b>🚪 Аудитория:</b> %d", *sub.LabAuditorium))
	} else if sub.LabDomain != nil {
		sb.WriteString(fmt.Sprintf("<b>⚛️ %s</b>", sub.LabDomain))
	}
	sb.WriteString(repeatLineBreaks(2))

	if sub.Weekday != nil {
		sb.WriteString(fmt.Sprintf("<b>📅 День:</b> %s", utils.WeekdayLocale[*sub.Weekday]))
		sb.WriteString(repeatLineBreaks(2))
	}

	if len(sub.PreferredTimes) > 0 {
		sb.WriteString(fmt.Sprintf("<b>🕐 Время:</b>"))
		sb.WriteString(repeatLineBreaks(2))
		for _, prefTime := range sub.PreferredTimes {
			sb.WriteString(fmt.Sprintf("<b>%s</b>", utils.TimeStartToLongLessonTime[prefTime.TimeStart]))
			sb.WriteString(repeatLineBreaks(1))
		}
	}
	return sb.String()
}

func UnsubSuccessMsg() string {
	return "<b>✅ Вы больше не подписаны на эту лабу</b>"
}

// ==

func NotifyMsg(notif *notification.Notification) string {
	slot := &notif.Slot
	var sb strings.Builder
	sb.WriteString("<b>🔥 Появилась запись!</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString(fmt.Sprintf("<b>⚛️ %s</b>", slot.Domain))
	sb.WriteString(repeatLineBreaks(2))
	longName := slot.Name
	if slot.Order != nil {
		longName += fmt.Sprintf(" (%d-ое место)", *slot.Order)
	}
	sb.WriteString(fmt.Sprintf("<b>📚 Лаба №%d. %s</b>", slot.Number, longName))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>🚪 Аудитория №%d</b>", slot.Auditorium))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🗓️ Когда:</b>")
	sb.WriteString(repeatLineBreaks(1))
	slotTimes := make([]time.Time, 0)
	for t := range slot.TimesTeachers {
		slotTimes = append(slotTimes, t)
	}
	grouped := utils.GroupTimesByDate(slotTimes)
	sortedDates := make([]time.Time, 0, len(grouped))
	for date := range grouped {
		sortedDates = append(sortedDates, date)
	}
	slices.SortFunc(sortedDates, func(a, b time.Time) int {
		return a.Compare(b)
	})
	for _, date := range sortedDates {
		dateRelative := utils.FormatDateRelative(date, time.Now())
		sb.WriteString(fmt.Sprintf("<b>⠀⠀%s:</b>", dateRelative))
		sb.WriteString(repeatLineBreaks(1))
		times := grouped[date]
		slices.SortFunc(times, func(a, b time.Time) int {
			return a.Compare(b)
		})
		for _, t := range times {
			stringParts := make([]string, 0)
			timeStart := t.Format("15:04")
			lessonTime := utils.TimeStartToShortLessonTime[timeStart]
			stringParts = append(stringParts, lessonTime)
			if teachers, ok := slot.TimesTeachers[t]; ok {
				stringParts = append(stringParts, teachers...)
			}
			if utils.IsTimeInPreferredTimes(t, &notif.PreferredTimes) {
				stringParts = append(stringParts, "⭐️ Ваше время")
			}
			sb.WriteString(fmt.Sprintf("<b>⠀⠀%s</b>", strings.Join(stringParts, " ")))
			sb.WriteString(repeatLineBreaks(1))
		}
		sb.WriteString(repeatLineBreaks(1))
	}
	return sb.String()
}

func repeatLineBreaks(breaks int) string {
	var sb strings.Builder
	for range breaks {
		sb.WriteString("\n")
	}
	return sb.String()
}

package cmd

import (
	"fmt"

	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/go-telegram/bot/models"
)

func selectLabDomainKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Механика", CallbackData: "lab_creation:domain:mechanics"}},
			{{Text: "Виртуалка", CallbackData: "lab_creation:domain:virtual"}},
			{{Text: "Электричество", CallbackData: "lab_creation:domain:electricity"}},
		},
	}
}

func selectLabWeekdayKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Понедельник", CallbackData: "lab_creation:weekday:1"}},
			{{Text: "Вторник", CallbackData: "lab_creation:weekday:2"}},
			{{Text: "Среда", CallbackData: "lab_creation:weekday:3"}},
			{{Text: "Четверг", CallbackData: "lab_creation:weekday:4"}},
			{{Text: "Пятница", CallbackData: "lab_creation:weekday:5"}},
			{{Text: "Суббота", CallbackData: "lab_creation:weekday:6"}},
			{{Text: "Воскресенье", CallbackData: "lab_creation:weekday:0"}},
			{{Text: "⏭️ Пропустить", CallbackData: "lab_creation:weekday:skip"}},
		},
	}
}

type Lesson struct {
	Text         string
	CallbackData string
}

var defaultLessons = []Lesson{
	{Text: "08:50 - 10:20 - 1️⃣ пара", CallbackData: "lab_creation:lesson:1"},
	{Text: "10:35 - 12:05 - 2️⃣ пара", CallbackData: "lab_creation:lesson:2"},
	{Text: "12:35 - 14:05 - 3️⃣ пара", CallbackData: "lab_creation:lesson:3"},
	{Text: "14:15 - 15:45 - 4️⃣ пара", CallbackData: "lab_creation:lesson:4"},
	{Text: "15:55 - 17:20 - 5️⃣ пара", CallbackData: "lab_creation:lesson:5"},
	{Text: "17:30 - 19:00 - 6️⃣ пара", CallbackData: "lab_creation:lesson:6"},
	{Text: "19:10 - 20:30 - 7️⃣ пара", CallbackData: "lab_creation:lesson:7"},
	{Text: "20:40 - 22:00 - 8️⃣ пара", CallbackData: "lab_creation:lesson:8"},
}

func selectLessonKbd(lessons []Lesson) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: make([][]models.InlineKeyboardButton, len(lessons)),
	}
	for idx, lesson := range lessons {
		keyboard.InlineKeyboard[idx] = []models.InlineKeyboardButton{
			{Text: lesson.Text, CallbackData: lesson.CallbackData},
		}
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
		{Text: "✅ Готово", CallbackData: "lab_creation:lesson:skip"},
	})

	return keyboard
}

func askLabConfirmationKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Создать", CallbackData: "lab_creation:confirm:create"},
				{Text: "❌ Отменить", CallbackData: "lab_creation:confirm:cancel"},
			},
		},
	}
}

func createUnsubKeyboard(subs []model.Subscription) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{},
	}
	for _, sub := range subs {
		label := fmt.Sprintf("Лаба №%d, ауд. №%d", sub.LabNumber, sub.LabAuditorium)
		var timeString string
		if sub.Weekday != nil && sub.DayTime != nil {
			timeString = fmt.Sprintf(", %s %s", weekDayLocale[int(*sub.Weekday)], timeLessonMap[*sub.DayTime])
		} else {
			timeString = ", Любое время"
		}
		label += timeString
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			[]models.InlineKeyboardButton{
				{Text: label, CallbackData: fmt.Sprintf("unsub:delete:%s", sub.UUID)},
			})
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		[]models.InlineKeyboardButton{
			{Text: "❌ Удалить все", CallbackData: "unsub:all"},
		})
	return keyboard
}

func createDeleteAllConfirmKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Да, удалить", CallbackData: "unsub:all:confirm"},
				{Text: "❌ Нет", CallbackData: "unsub:all:cancel"},
			},
		},
	}
}

func createLinkKeyboard(url string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🔗 ЗАПИСАТЬСЯ", URL: url},
			},
		},
	}
}

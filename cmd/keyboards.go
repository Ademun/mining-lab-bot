package cmd

import (
	"fmt"

	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/go-telegram/bot/models"
)

func createWeekdayKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Понедельник", CallbackData: "weekday:1"}},
			{{Text: "Вторник", CallbackData: "weekday:2"}},
			{{Text: "Среда", CallbackData: "weekday:3"}},
			{{Text: "Четверг", CallbackData: "weekday:4"}},
			{{Text: "Пятница", CallbackData: "weekday:5"}},
			{{Text: "Суббота", CallbackData: "weekday:6"}},
			{{Text: "Воскресенье", CallbackData: "weekday:0"}},
			{{Text: "⏭️ Пропустить", CallbackData: "skip:weekday"}},
		},
	}
}

func createLessonKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "08:50 - 10:20 - 1️⃣ пара", CallbackData: "lesson:1"}},
			{{Text: "10:35 - 12:05 - 2️⃣ пара", CallbackData: "lesson:2"}},
			{{Text: "12:35 - 14:05 - 3️⃣ пара", CallbackData: "lesson:3"}},
			{{Text: "14:15 - 15:45 - 4️⃣ пара", CallbackData: "lesson:4"}},
			{{Text: "15:55 - 17:20 - 5️⃣ пара", CallbackData: "lesson:5"}},
			{{Text: "17:30 - 19:00 - 6️⃣ пара", CallbackData: "lesson:6"}},
			{{Text: "19:10 - 20:30 - 7️⃣ пара", CallbackData: "lesson:7"}},
			{{Text: "20:40 - 22:00 - 8️⃣ пара", CallbackData: "lesson:8"}},
		},
	}
}

func createConfirmationKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Создать", CallbackData: "confirm:create"},
				{Text: "❌ Отменить", CallbackData: "confirm:cancel"},
			},
		},
	}
}

func createUnsubKeyboard(subs []model.Subscription) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{},
	}
	for idx, sub := range subs {
		label := fmt.Sprintf("%d. Лаба №%d, ауд. №%d", idx+1, sub.LabNumber,
			sub.LabAuditorium)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			[]models.InlineKeyboardButton{
				{Text: label, CallbackData: fmt.Sprintf("unsub:view:%s", sub.UUID)},
				{Text: "❌", CallbackData: fmt.Sprintf("unsub:delete:%s", sub.UUID)},
			})
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		[]models.InlineKeyboardButton{
			{Text: "🗑️ Удалить все", CallbackData: "unsub:all"},
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

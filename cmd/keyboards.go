package cmd

import (
	"fmt"

	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/go-telegram/bot/models"
)

func createWeekdayKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "ПН", CallbackData: "weekday:1"},
				{Text: "ВТ", CallbackData: "weekday:2"},
				{Text: "СР", CallbackData: "weekday:3"},
				{Text: "ЧТ", CallbackData: "weekday:4"},
			},
			{
				{Text: "ПТ", CallbackData: "weekday:5"},
				{Text: "СБ", CallbackData: "weekday:6"},
				{Text: "ВС", CallbackData: "weekday:0"},
			},
			{
				{Text: "⏭️ Пропустить", CallbackData: "skip:weekday"},
			},
		},
	}
}

func createLessonKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "1️⃣ пара 8:50-10:20", CallbackData: "lesson:1"}},
			{{Text: "2️⃣ пара 10:35-12:05", CallbackData: "lesson:2"}},
			{{Text: "3️⃣ пара 12:35-14:05", CallbackData: "lesson:3"}},
			{{Text: "4️⃣ пара 14:15-15:45", CallbackData: "lesson:4"}},
			{{Text: "5️⃣ пара 15:55-17:20", CallbackData: "lesson:5"}},
			{{Text: "6️⃣ пара 17:30-19:00", CallbackData: "lesson:6"}},
			{{Text: "7️⃣ пара 19:10-20:30", CallbackData: "lesson:7"}},
			{{Text: "8️⃣ пара 20:40-22:00", CallbackData: "lesson:8"}},
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

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

func createSkipKeyboard(field string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "⏭️ Пропустить", CallbackData: fmt.Sprintf("skip:%s", field)},
			},
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

package cmd

import (
	"fmt"

	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

func selectLabTypeKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Выполнение", CallbackData: "sub_creation:type:performance"}},
			{{Text: "Защита", CallbackData: "sub_creation:domain:virtual:defence"}},
		},
	}
}

func selectLabDomainKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Механика", CallbackData: "sub_creation:domain:mechanics"}},
			{{Text: "Виртуалка", CallbackData: "sub_creation:domain:virtual"}},
			{{Text: "Электричество", CallbackData: "sub_creation:domain:electricity"}},
		},
	}
}

func selectLabWeekdayKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Понедельник", CallbackData: "sub_creation:weekday:1"}},
			{{Text: "Вторник", CallbackData: "sub_creation:weekday:2"}},
			{{Text: "Среда", CallbackData: "sub_creation:weekday:3"}},
			{{Text: "Четверг", CallbackData: "sub_creation:weekday:4"}},
			{{Text: "Пятница", CallbackData: "sub_creation:weekday:5"}},
			{{Text: "Суббота", CallbackData: "sub_creation:weekday:6"}},
			{{Text: "Воскресенье", CallbackData: "sub_creation:weekday:0"}},
			{{Text: "⏭️ Пропустить", CallbackData: "sub_creation:weekday:skip"}},
		},
	}
}

type Lesson struct {
	Text         string
	CallbackData string
}

var defaultLessons = []Lesson{
	{Text: "08:50 - 10:20 - 1️⃣ пара", CallbackData: "sub_creation:lesson:1"},
	{Text: "10:35 - 12:05 - 2️⃣ пара", CallbackData: "sub_creation:lesson:2"},
	{Text: "12:35 - 14:05 - 3️⃣ пара", CallbackData: "sub_creation:lesson:3"},
	{Text: "14:15 - 15:45 - 4️⃣ пара", CallbackData: "sub_creation:lesson:4"},
	{Text: "15:55 - 17:20 - 5️⃣ пара", CallbackData: "sub_creation:lesson:5"},
	{Text: "17:30 - 19:00 - 6️⃣ пара", CallbackData: "sub_creation:lesson:6"},
	{Text: "19:10 - 20:30 - 7️⃣ пара", CallbackData: "sub_creation:lesson:7"},
	{Text: "20:40 - 22:00 - 8️⃣ пара", CallbackData: "sub_creation:lesson:8"},
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
		{Text: "✅ Готово", CallbackData: "sub_creation:lesson:skip"},
	})

	return keyboard
}

func askLabConfirmationKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Создать", CallbackData: "sub_creation:confirm:create"},
				{Text: "❌ Отменить", CallbackData: "sub_creation:confirm:cancel"},
			},
		},
	}
}

func listSubsKbd(subUUID uuid.UUID, subIdx, totalSubs int) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: make([][]models.InlineKeyboardButton, 0),
	}
	paginationRow := make([]models.InlineKeyboardButton, 0)
	if subIdx > 0 {
		paginationRow = append(paginationRow, models.InlineKeyboardButton{
			Text: "<<", CallbackData: fmt.Sprintf("sub_list:move:%d", subIdx-1),
		})
	}
	paginationRow = append(paginationRow, models.InlineKeyboardButton{
		Text: fmt.Sprintf("%d/%d", subIdx+1, totalSubs), CallbackData: fmt.Sprintf("sub_list:move:%d", subIdx),
	})
	if subIdx < totalSubs-1 {
		paginationRow = append(paginationRow, models.InlineKeyboardButton{
			Text: ">>", CallbackData: fmt.Sprintf("sub_list:move:%d", subIdx+1),
		})
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, paginationRow)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
		{
			Text: "🗑️ Удалить", CallbackData: fmt.Sprintf("sub_list:delete:%s", subUUID.String()),
		},
	})
	return keyboard
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

package presentation

import (
	"fmt"

	"github.com/Ademun/mining-lab-bot/cmd/internal/utils"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

// Универсальные клавиатуры

func CancelKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "❌ Отменить", CallbackData: "cancel"}},
		},
	}
}

func SelectWeekdayKbd(withSkip bool) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Понедельник", CallbackData: "weekday:1"}},
			{{Text: "Вторник", CallbackData: "weekday:2"}},
			{{Text: "Среда", CallbackData: "weekday:3"}},
			{{Text: "Четверг", CallbackData: "weekday:4"}},
			{{Text: "Пятница", CallbackData: "weekday:5"}},
			{{Text: "Суббота", CallbackData: "weekday:6"}},
			{{Text: "Воскресенье", CallbackData: "weekday:0"}},
		},
	}

	if withSkip {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
			{Text: "⏭️ Пропустить", CallbackData: "weekday:skip"},
		})
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
		{Text: "❌ Отменить", CallbackData: "cancel"},
	})

	return keyboard
}

func SelectLessonKbd(lessons []utils.Lesson, multi bool) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: make([][]models.InlineKeyboardButton, len(lessons)),
	}
	for idx, lesson := range lessons {
		keyboard.InlineKeyboard[idx] = []models.InlineKeyboardButton{
			{Text: lesson.Text, CallbackData: lesson.CallbackData},
		}
	}

	if multi {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, [][]models.InlineKeyboardButton{
			{{Text: "✅ Готово", CallbackData: "lesson:skip"}},
		}...)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, [][]models.InlineKeyboardButton{
		{{Text: "❌ Отменить", CallbackData: "cancel"}},
	}...)

	return keyboard
}

// Subscription creation keyboards

func SelectLabTypeKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Выполнение", CallbackData: "type:performance"}},
			{{Text: "Защита", CallbackData: "type:defence"}},
			{{Text: "❌ Отменить создание", CallbackData: "cancel"}},
		},
	}
}

func SelectLabDomainKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Механика", CallbackData: "domain:mechanics"}},
			{{Text: "Виртуалка", CallbackData: "domain:virtual"}},
			{{Text: "Электричество", CallbackData: "domain:electricity"}},
			{{Text: "❌ Отменить создание", CallbackData: "cancel"}},
		},
	}
}

func AskSubCreationConfirmationKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Создать", CallbackData: "confirm:create"},
				{Text: "❌ Отменить", CallbackData: "cancel"},
			},
		},
	}
}

// Subscription listing keyboards

func ListSubsKbd(subUUID uuid.UUID, subIdx, totalSubs int) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: make([][]models.InlineKeyboardButton, 0),
	}
	paginationRow := make([]models.InlineKeyboardButton, 0)
	if subIdx > 0 {
		paginationRow = append(paginationRow, models.InlineKeyboardButton{
			Text: "<<", CallbackData: fmt.Sprintf("move:%d", subIdx-1),
		})
	}
	paginationRow = append(paginationRow, models.InlineKeyboardButton{
		Text: fmt.Sprintf("%d/%d", subIdx+1, totalSubs), CallbackData: fmt.Sprintf("move:%d", subIdx),
	})
	if subIdx < totalSubs-1 {
		paginationRow = append(paginationRow, models.InlineKeyboardButton{
			Text: ">>", CallbackData: fmt.Sprintf("move:%d", subIdx+1),
		})
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, paginationRow)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
		{
			Text: "🗑️ Удалить", CallbackData: fmt.Sprintf("delete:%s", subUUID.String()),
		},
	})
	return keyboard
}

func LinkKbd(url string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🔗 ЗАПИСАТЬСЯ", URL: url},
			},
		},
	}
}

// Teacher report keyboards

func SelectWeekParityKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Чётная", CallbackData: "parity:even"}},
			{{Text: "Нечётная", CallbackData: "parity:odd"}},
			{{Text: "❌ Отменить", CallbackData: "cancel"}},
		},
	}
}

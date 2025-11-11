package presentation

import (
	"fmt"

	"github.com/Ademun/mining-lab-bot/cmd/internal/utils"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

func SelectLabTypeKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Выполнение", CallbackData: "sub_creation:type:performance"}},
			{{Text: "Защита", CallbackData: "sub_creation:type:defence"}},
		},
	}
}

func SelectLabDomainKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "Механика", CallbackData: "sub_creation:domain:mechanics"}},
			{{Text: "Виртуалка", CallbackData: "sub_creation:domain:virtual"}},
			{{Text: "Электричество", CallbackData: "sub_creation:domain:electricity"}},
		},
	}
}

func SelectWeekdayKbd() *models.InlineKeyboardMarkup {
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

func SelectLessonKbd(lessons []utils.Lesson) *models.InlineKeyboardMarkup {
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

func AskSubCreationConfirmationKbd() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Создать", CallbackData: "sub_creation:confirm:create"},
				{Text: "❌ Отменить", CallbackData: "sub_creation:confirm:cancel"},
			},
		},
	}
}

func ListSubsKbd(subUUID uuid.UUID, subIdx, totalSubs int) *models.InlineKeyboardMarkup {
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

func LinkKbd(url string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🔗 ЗАПИСАТЬСЯ", CallbackData: url},
			},
		},
	}
}

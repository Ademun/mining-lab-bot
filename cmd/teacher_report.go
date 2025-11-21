package cmd

import (
	"context"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/cmd/internal/presentation"
	"github.com/Ademun/mining-lab-bot/cmd/internal/utils"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *telegramBot) handleTeacherReport(ctx context.Context, api *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID

	newData := &fsm.TeacherReportFlowData{
		UserID: userID,
	}

	b.TryTransition(ctx, userID, fsm.StepAwaitingTeacherAuditorium, newData)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.AskLabAuditoriumMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.CancelKbd(),
	})
}

func (b *telegramBot) handleTeacherAuditorium(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleTeacherReportCancellation(ctx, b, update) {
		return
	}
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID
	auditoriumStr := update.Message.Text

	auditorium, cause := validateLabAuditorium(auditoriumStr)
	if cause != "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.ValidationErrorMsg(cause),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	newData, ok := data.(*fsm.TeacherReportFlowData)
	if !ok {
		slog.Error("Critical error: unable to assert flow data",
			"data", data,
			"service", logger.TelegramBot)
		b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.GenericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	newData.Auditorium = auditorium

	b.TryTransition(ctx, userID, fsm.StepAwaitingTeacherWeekParity, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.AskWeekParityMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.SelectWeekParityKbd(),
	})
}

func (b *telegramBot) handleTeacherWeekParity(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleTeacherReportCancellation(ctx, b, update) {
		return
	}
	if update.CallbackQuery == nil {
		return
	}
	userID := update.CallbackQuery.From.ID
	weekParity := extractWeekParity(update)

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	newData, ok := data.(*fsm.TeacherReportFlowData)
	if !ok {
		slog.Error("Critical error: unable to assert flow data",
			"data", data,
			"service", logger.TelegramBot)
		b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.GenericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	newData.WeekParity = weekParity

	b.TryTransition(ctx, userID, fsm.StepAwaitingTeacherWeekday, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.AskTeacherWeekdayMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.SelectWeekdayKbd(false),
	})
}

func (b *telegramBot) handleTeacherWeekday(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleTeacherReportCancellation(ctx, b, update) {
		return
	}
	if update.CallbackQuery == nil {
		return
	}
	userID := update.CallbackQuery.From.ID
	weekday := extractWeekday(update)

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	newData, ok := data.(*fsm.TeacherReportFlowData)
	if !ok {
		slog.Error("Critical error: unable to assert flow data",
			"data", data,
			"service", logger.TelegramBot)
		b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.GenericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	if weekday != nil {
		newData.Weekday = *weekday
	}

	b.TryTransition(ctx, userID, fsm.StepAwaitingTeacherLesson, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.AskTeacherLessonMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.SelectLessonKbd(utils.DefaultLessons, false),
	})
}

func (b *telegramBot) handleTeacherLesson(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleTeacherReportCancellation(ctx, b, update) {
		return
	}
	if update.CallbackQuery == nil {
		return
	}
	userID := update.CallbackQuery.From.ID
	lessonNum := extractLesson(update)

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	newData, ok := data.(*fsm.TeacherReportFlowData)
	if !ok {
		slog.Error("Critical error: unable to assert flow data",
			"data", data,
			"service", logger.TelegramBot)
		b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.GenericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if lessonNum != nil {
		newData.LessonNum = *lessonNum
	}

	b.TryTransition(ctx, userID, fsm.StepAwaitingTeacherSurname, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.AskTeacherSurnameMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.CancelKbd(),
	})
}

func (b *telegramBot) handleTeacherSurname(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleTeacherReportCancellation(ctx, b, update) {
		return
	}
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID
	surnameStr := update.Message.Text

	surname, cause := validateTeacherSurname(surnameStr)
	if cause != "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.ValidationErrorMsg(cause),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	newData, ok := data.(*fsm.TeacherReportFlowData)
	if !ok {
		slog.Error("Critical error: unable to assert flow data",
			"data", data,
			"service", logger.TelegramBot)
		b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.GenericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	newData.Surname = surname

	b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      presentation.TeacherReportSuccessMsg(),
		ParseMode: models.ParseModeHTML,
	})

	// Отправка данных админу
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: b.options.AdminID,
		Text: presentation.TeacherReportAdminMsg(
			newData.UserID,
			newData.Auditorium,
			newData.WeekParity,
			newData.Weekday,
			newData.LessonNum,
			newData.Surname,
		),
		ParseMode: models.ParseModeHTML,
	})
}

func handleTeacherReportCancellation(ctx context.Context, b *telegramBot, update *models.Update) bool {
	if update.CallbackQuery == nil {
		return false
	}
	userID := update.CallbackQuery.From.ID
	cancelledStr := update.CallbackQuery.Data
	if cancelledStr != "cancel" {
		return false
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	b.TryTransition(ctx, userID, fsm.StepIdle, nil)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      presentation.TeacherReportCancelledMsg(),
		ParseMode: models.ParseModeHTML,
	})

	return true
}

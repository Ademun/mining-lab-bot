package cmd

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/cmd/internal/presentation"
	"github.com/Ademun/mining-lab-bot/cmd/internal/utils"
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *telegramBot) handleSubscriptionCreation(ctx context.Context, api *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID

	newData := &fsm.SubscriptionCreationFlowData{
		UserID: int(userID),
	}

	b.TryTransition(ctx, userID, fsm.StepAwaitingLabType, newData)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.AskLabTypeMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.SelectLabTypeKbd(),
	})
}

func (b *telegramBot) handleLabType(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleCancellation(ctx, b, update) {
		return
	}
	if update.CallbackQuery == nil {
		return
	}
	userID := update.CallbackQuery.From.ID
	labType := extractLabType(update)

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
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
	newData.LabType = labType

	b.TryTransition(ctx, userID, fsm.StepAwaitingLabNumber, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.AskLabNumberMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.SubCreationCancelKbd(),
	})
}

func (b *telegramBot) handleLabNumber(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleCancellation(ctx, b, update) {
		return
	}
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID
	labNumberStr := update.Message.Text

	labNumber, cause := validateLabNumber(labNumberStr)
	if cause != "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.ValidationErrorMsg(cause),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
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
	newData.LabNumber = labNumber

	switch newData.LabType {
	case polling.LabTypePerformance:
		b.TryTransition(ctx, userID, fsm.StepAwaitingLabAuditorium, newData)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      userID,
			Text:        presentation.AskLabAuditoriumMsg(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: presentation.SubCreationCancelKbd(),
		})
	case polling.LabTypeDefence:
		b.TryTransition(ctx, userID, fsm.StepAwaitingLabDomain, newData)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      userID,
			Text:        presentation.AskLabDomainMsg(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: presentation.SelectLabDomainKbd(),
		})
	}
}

func (b *telegramBot) handleLabAuditorium(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleCancellation(ctx, b, update) {
		return
	}
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID
	labAuditoriumStr := update.Message.Text

	labAuditorium, cause := validateLabAuditorium(labAuditoriumStr)
	if cause != "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.ValidationErrorMsg(cause),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
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
	}
	newData.LabAuditorium = &labAuditorium

	b.TryTransition(ctx, userID, fsm.StepAwaitingLabWeekday, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.AskWeekdayMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.SelectWeekdayKbd(),
	})
}

func (b *telegramBot) handleLabDomain(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleCancellation(ctx, b, update) {
		return
	}
	if update.CallbackQuery == nil {
		return
	}
	userID := update.CallbackQuery.From.ID
	labDomain := extractLabDomain(update)

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
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
	newData.LabDomain = labDomain

	b.TryTransition(ctx, userID, fsm.StepAwaitingLabWeekday, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.AskWeekdayMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.SelectWeekdayKbd(),
	})
}

func (b *telegramBot) handleWeekday(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleCancellation(ctx, b, update) {
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

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
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
	newData.Weekday = weekday

	if weekday == nil {
		b.TryTransition(ctx, userID, fsm.StepAwaitingSubCreationConfirmation, newData)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      userID,
			Text:        presentation.AskSubCreationConfirmationMsg(parseFlowData(newData)),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: presentation.AskSubCreationConfirmationKbd(),
		})
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		return
	}

	b.TryTransition(ctx, userID, fsm.StepAwaitingLabLessons, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.AskLessonsMsg(nil),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.SelectLessonKbd(utils.DefaultLessons),
	})
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
}

func (b *telegramBot) handleLessons(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleCancellation(ctx, b, update) {
		return
	}
	if update.CallbackQuery == nil {
		return
	}
	userID := update.CallbackQuery.From.ID
	lesson := extractLesson(update)

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
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

	if lesson == nil {
		sub := parseFlowData(newData)
		b.TryTransition(ctx, userID, fsm.StepAwaitingSubCreationConfirmation, newData)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      userID,
			Text:        presentation.AskSubCreationConfirmationMsg(sub),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: presentation.AskSubCreationConfirmationKbd(),
		})
		return
	}

	newData.Lessons = append(newData.Lessons, *lesson)

	existingLessonsMap := make(map[int]bool)
	for _, lesson := range newData.Lessons {
		existingLessonsMap[lesson] = true
	}

	kbdLessons := make([]utils.Lesson, 0, len(utils.DefaultLessons))
	for i, lesson := range utils.DefaultLessons {
		lessonNum := i + 1
		if !existingLessonsMap[lessonNum] {
			kbdLessons = append(kbdLessons, lesson)
		}
	}

	b.TryTransition(ctx, userID, fsm.StepAwaitingLabLessons, newData)
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    userID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      presentation.AskLessonsMsg(newData.Lessons),
		ParseMode: models.ParseModeHTML,
	})
	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      userID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: presentation.SelectLessonKbd(kbdLessons),
	})
}

func (b *telegramBot) handleSubCreationConfirmation(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if handleCancellation(ctx, b, update) {
		return
	}
	if update.CallbackQuery == nil {
		return
	}
	userID := update.CallbackQuery.From.ID

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
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

	sub := parseFlowData(newData)
	b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})

	err := b.subscriptionService.Subscribe(ctx, *sub)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.GenericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      presentation.SubCreationSuccessMsg(),
		ParseMode: models.ParseModeHTML,
	})
	b.notifService.NotifyNewSubscription(ctx, *sub)
	return
}

func extractLabType(update *models.Update) polling.LabType {
	labTypeStr := update.CallbackQuery.Data
	labTypeStr = strings.TrimPrefix(labTypeStr, "sub_creation:type:")

	var labType polling.LabType
	switch labTypeStr {
	case "performance":
		labType = polling.LabTypePerformance
	case "defence":
		labType = polling.LabTypeDefence
	}
	return labType
}

func extractLabDomain(update *models.Update) *polling.LabDomain {
	labDomainStr := update.CallbackQuery.Data
	labDomainStr = strings.TrimPrefix(labDomainStr, "sub_creation:domain:")

	var labDomain polling.LabDomain
	switch labDomainStr {
	case "mechanics":
		labDomain = polling.LabDomainMechanics
	case "virtual":
		labDomain = polling.LabDomainVirtual
	case "electricity":
		labDomain = polling.LabDomainElectricity
	}
	return &labDomain
}

func extractWeekday(update *models.Update) *int {
	labWeekdayStr := update.CallbackQuery.Data
	labWeekdayStr = strings.TrimPrefix(labWeekdayStr, "sub_creation:weekday:")

	if labWeekdayStr == "skip" {
		return nil
	}
	labWeekdayInt, _ := strconv.Atoi(labWeekdayStr)
	return &labWeekdayInt
}

func extractLesson(update *models.Update) *int {
	labLessonStr := update.CallbackQuery.Data
	labLessonStr = strings.TrimPrefix(labLessonStr, "sub_creation:lesson:")

	if labLessonStr == "skip" {
		return nil
	}
	labLessonInt, _ := strconv.Atoi(labLessonStr)
	return &labLessonInt
}

func handleCancellation(ctx context.Context, b *telegramBot, update *models.Update) bool {
	if update.CallbackQuery == nil {
		return false
	}
	userID := update.CallbackQuery.From.ID
	cancelledStr := update.CallbackQuery.Data
	if strings.TrimPrefix(cancelledStr, "sub_creation:") != "cancel" {
		return false
	}

	b.TryTransition(ctx, userID, fsm.StepIdle, nil)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      presentation.SubCreationCancelledMsg(),
		ParseMode: models.ParseModeHTML,
	})

	return true
}

func parseFlowData(data *fsm.SubscriptionCreationFlowData) *subscription.RequestSubscription {
	return &subscription.RequestSubscription{
		UserID:        data.UserID,
		Type:          data.LabType,
		LabNumber:     data.LabNumber,
		LabAuditorium: data.LabAuditorium,
		LabDomain:     data.LabDomain,
		Weekday:       data.Weekday,
		Lessons:       data.Lessons,
	}
}

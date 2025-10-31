package cmd

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *telegramBot) handleLabType(ctx context.Context, api *bot.Bot, update *models.Update, state *fsm.State) {
	userID := update.CallbackQuery.From.ID
	labType := extractLabType(update)

	slog.Info("Handling lab type selection",
		"user_id", userID,
		"lab_type", labType,
		"service", logger.TelegramBot)

	newData := map[string]interface{}{
		"user_id":  userID,
		"lab_type": labType,
	}
	b.tryTransition(ctx, api, userID, fsm.StepAwaitingLabNumber, newData)

	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      askLabNumberMsg(),
		ParseMode: models.ParseModeHTML,
	}); err != nil {
		slog.Error("Failed to send lab number request",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}

	if _, err := api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	}); err != nil {
		slog.Error("Failed to answer callback query",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

func (b *telegramBot) handleLabNumber(ctx context.Context, api *bot.Bot, update *models.Update, state *fsm.State) {
	userID := update.Message.From.ID
	labNumberStr := update.Message.Text

	slog.Debug("Handling lab number input",
		"user_id", userID,
		"input", labNumberStr,
		"service", logger.TelegramBot)

	labNumber, err := strconv.Atoi(labNumberStr)
	if err != nil {
		slog.Warn("Invalid lab number format",
			"error", err,
			"user_id", userID,
			"input", labNumberStr,
			"service", logger.TelegramBot)

		if _, sendErr := api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   labNumberValidationErrorMsg(),
		}); sendErr != nil {
			slog.Error("Failed to send validation error message",
				"error", sendErr,
				"user_id", userID,
				"service", logger.TelegramBot)
		}
		return
	}

	slog.Info("Lab number validated",
		"user_id", userID,
		"lab_number", labNumber,
		"service", logger.TelegramBot)

	newData := map[string]interface{}{
		"lab_number": labNumber,
	}

	labType, ok := state.Data["lab_type"].(polling.LabType)
	if !ok {
		slog.Error("Failed to get lab type from state",
			"user_id", userID,
			"service", logger.TelegramBot)
		return
	}

	switch labType {
	case polling.LabTypePerformance:
		slog.Debug("Processing performance lab type",
			"user_id", userID,
			"service", logger.TelegramBot)

		b.tryTransition(ctx, api, userID, fsm.StepAwaitingLabAuditorium, newData)
		if _, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      askLabAuditoriumMsg(),
			ParseMode: models.ParseModeHTML,
		}); err != nil {
			slog.Error("Failed to send auditorium request",
				"error", err,
				"user_id", userID,
				"service", logger.TelegramBot)
		}

	case polling.LabTypeDefence:
		slog.Debug("Processing defence lab type",
			"user_id", userID,
			"service", logger.TelegramBot)

		b.tryTransition(ctx, api, userID, fsm.StepAwaitingLabDomain, newData)
		if _, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      userID,
			Text:        askLabDomainMsg(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: selectLabDomainKbd(),
		}); err != nil {
			slog.Error("Failed to send domain request",
				"error", err,
				"user_id", userID,
				"service", logger.TelegramBot)
		}
	}
}

func (b *telegramBot) handleLabAuditorium(ctx context.Context, api *bot.Bot, update *models.Update, state *fsm.State) {
	userID := update.Message.From.ID
	labAuditoriumStr := update.Message.Text

	slog.Debug("Handling lab auditorium input",
		"user_id", userID,
		"input", labAuditoriumStr,
		"service", logger.TelegramBot)

	labAuditorium, err := strconv.Atoi(labAuditoriumStr)
	if err != nil {
		slog.Warn("Invalid lab auditorium format",
			"error", err,
			"user_id", userID,
			"input", labAuditoriumStr,
			"service", logger.TelegramBot)

		if _, sendErr := api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      labAuditoriumValidationErrorMsg(),
			ParseMode: models.ParseModeHTML,
		}); sendErr != nil {
			slog.Error("Failed to send validation error message",
				"error", sendErr,
				"user_id", userID,
				"service", logger.TelegramBot)
		}
		return
	}

	slog.Info("Lab auditorium validated",
		"user_id", userID,
		"lab_auditorium", labAuditorium,
		"service", logger.TelegramBot)

	newData := map[string]interface{}{
		"lab_auditorium": labAuditorium,
	}

	b.tryTransition(ctx, api, userID, fsm.StepAwaitingLabWeekday, newData)
	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        askLabWeekdayMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: selectLabWeekdayKbd(),
	}); err != nil {
		slog.Error("Failed to send weekday request",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

func (b *telegramBot) handleLabDomain(ctx context.Context, api *bot.Bot, update *models.Update, state *fsm.State) {
	userID := update.CallbackQuery.From.ID
	labDomain := extractLabDomain(update)

	slog.Info("Handling lab domain selection",
		"user_id", userID,
		"lab_domain", labDomain,
		"service", logger.TelegramBot)

	newData := map[string]interface{}{
		"lab_domain": labDomain,
	}

	b.tryTransition(ctx, api, userID, fsm.StepAwaitingLabWeekday, newData)
	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        askLabWeekdayMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: selectLabWeekdayKbd(),
	}); err != nil {
		slog.Error("Failed to send weekday request",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}

	if _, err := api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	}); err != nil {
		slog.Error("Failed to answer callback query",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

func (b *telegramBot) handleLabWeekday(ctx context.Context, api *bot.Bot, update *models.Update, state *fsm.State) {
	userID := update.CallbackQuery.From.ID
	labWeekday := extractLabWeekday(update)

	slog.Info("Handling lab weekday selection",
		"user_id", userID,
		"lab_weekday", labWeekday,
		"service", logger.TelegramBot)

	newData := map[string]interface{}{
		"lab_weekday": labWeekday,
	}

	b.tryTransition(ctx, api, userID, fsm.StepAwaitingLabLessons, newData)
	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        askLabLessonMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: selectLessonKbd(defaultLessons),
	}); err != nil {
		slog.Error("Failed to send lesson request",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}

	if _, err := api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	}); err != nil {
		slog.Error("Failed to answer callback query",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

func (b *telegramBot) handleLabLessons(ctx context.Context, api *bot.Bot, update *models.Update, state *fsm.State) {
	userID := update.CallbackQuery.From.ID
	labLesson := extractLabLesson(update)

	slog.Debug("Handling lab lesson selection",
		"user_id", userID,
		"lab_lesson", labLesson,
		"service", logger.TelegramBot)

	if labLesson == nil {
		slog.Info("User finished lesson selection, requesting confirmation",
			"user_id", userID,
			"service", logger.TelegramBot)

		sub := parseState(userID, state)
		b.tryTransition(ctx, api, userID, fsm.StepAwaitingLabConfirmation, nil)

		if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      userID,
			Text:        askLabConfirmationMsg(sub),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: askLabConfirmationKbd(),
		}); err != nil {
			slog.Error("Failed to send confirmation request",
				"error", err,
				"user_id", userID,
				"service", logger.TelegramBot)
		}
		return
	}

	prevLessons, _ := state.Data["lab_lessons"].([]int)
	newLessons := append(prevLessons, *labLesson)

	slog.Debug("Adding lesson to selection",
		"user_id", userID,
		"lesson", *labLesson,
		"total_lessons", len(newLessons),
		"service", logger.TelegramBot)

	newData := map[string]interface{}{
		"lab_lessons": newLessons,
	}

	// Исправленная логика фильтрации уроков
	existingLessonsMap := make(map[int]bool)
	for _, lesson := range newLessons {
		existingLessonsMap[lesson] = true
	}

	kbdLessons := make([]Lesson, 0, len(defaultLessons))
	for _, lesson := range defaultLessons {
		lessonNum, err := extractLessonNumberFromCallback(lesson.CallbackData)
		if err != nil {
			slog.Warn("Failed to parse lesson number from callback",
				"error", err,
				"callback_data", lesson.CallbackData,
				"service", logger.TelegramBot)
			continue
		}
		if !existingLessonsMap[lessonNum] {
			kbdLessons = append(kbdLessons, lesson)
		}
	}

	b.tryTransition(ctx, api, userID, fsm.StepAwaitingLabLessons, newData)

	if _, err := api.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      userID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: selectLessonKbd(kbdLessons),
	}); err != nil {
		slog.Error("Failed to update lesson keyboard",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

func (b *telegramBot) handleLabConfirmation(ctx context.Context, api *bot.Bot, update *models.Update, state *fsm.State) {
	userID := update.CallbackQuery.From.ID
	confirmed := extractConfirmed(update)

	slog.Info("Handling lab confirmation",
		"user_id", userID,
		"confirmed", confirmed,
		"service", logger.TelegramBot)

	if confirmed {
		sub := parseState(userID, state)
		b.tryTransition(ctx, api, userID, fsm.StepIdle, nil)

		err := b.subscriptionService.Subscribe(ctx, *sub)
		if err != nil {
			slog.Error("Failed to create subscription",
				"error", err,
				"user_id", userID,
				"service", logger.TelegramBot)

			if _, sendErr := api.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    userID,
				Text:      subCreationErrorMessage(err),
				ParseMode: models.ParseModeHTML,
			}); sendErr != nil {
				slog.Error("Failed to send error message",
					"error", sendErr,
					"user_id", userID,
					"service", logger.TelegramBot)
			}
			return
		}

		slog.Info("Subscription created successfully",
			"user_id", userID,
			"service", logger.TelegramBot)

		if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      subCreationSuccessMsg(),
			ParseMode: models.ParseModeHTML,
		}); err != nil {
			slog.Error("Failed to send success message",
				"error", err,
				"user_id", userID,
				"service", logger.TelegramBot)
		}
		return
	}

	slog.Info("Subscription creation cancelled by user",
		"user_id", userID,
		"service", logger.TelegramBot)

	b.tryTransition(ctx, api, userID, fsm.StepIdle, nil)

	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      subCreationCancelledMessage(),
		ParseMode: models.ParseModeHTML,
	}); err != nil {
		slog.Error("Failed to send cancellation message",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

func (b *telegramBot) tryTransition(ctx context.Context, api *bot.Bot, userID int64, newStep fsm.ConversationStep, newData map[string]interface{}) {
	if err := b.router.Transition(ctx, userID, newStep, newData); err != nil {
		slog.Error("State transition failed",
			"error", err,
			"user_id", userID,
			"new_step", newStep,
			"service", logger.TelegramBot)

		if _, sendErr := api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      genericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		}); sendErr != nil {
			slog.Error("Failed to send error message",
				"error", sendErr,
				"user_id", userID,
				"service", logger.TelegramBot)
		}
	}
}

func extractLabType(update *models.Update) polling.LabType {
	labTypeStr := update.CallbackQuery.Data
	labTypeStr = strings.TrimPrefix(labTypeStr, "lab_creation:type:")

	var labType polling.LabType
	switch labTypeStr {
	case "performance":
		labType = polling.LabTypePerformance
	case "defence":
		labType = polling.LabTypeDefence
	}
	return labType
}

func extractLabDomain(update *models.Update) polling.LabDomain {
	labDomainStr := update.CallbackQuery.Data
	labDomainStr = strings.TrimPrefix(labDomainStr, "lab_creation:domain:")

	var labDomain polling.LabDomain
	switch labDomainStr {
	case "mechanics":
		labDomain = polling.LabDomainMechanics
	case "virtual":
		labDomain = polling.LabDomainVirtual
	case "electricity":
		labDomain = polling.LabDomainElectricity
	}
	return labDomain
}

func extractLabWeekday(update *models.Update) *int {
	labWeekdayStr := update.CallbackQuery.Data
	labWeekdayStr = strings.TrimPrefix(labWeekdayStr, "lab_creation:weekday:")

	if labWeekdayStr == "skip" {
		return nil
	}
	labWeekdayInt, _ := strconv.Atoi(labWeekdayStr)
	return &labWeekdayInt
}

func extractLabLesson(update *models.Update) *int {
	labLessonStr := update.CallbackQuery.Data
	labLessonStr = strings.TrimPrefix(labLessonStr, "lab_creation:lesson:")

	if labLessonStr == "skip" {
		return nil
	}
	labLessonInt, _ := strconv.Atoi(labLessonStr)
	return &labLessonInt
}

func extractConfirmed(update *models.Update) bool {
	confirmedStr := update.CallbackQuery.Data
	confirmedStr = strings.TrimPrefix(confirmedStr, "lab_creation:confirm:")
	return confirmedStr == "create"
}

func extractLessonNumberFromCallback(callbackData string) (int, error) {
	lessonStr := strings.TrimPrefix(callbackData, "lab_creation:lesson:")
	return strconv.Atoi(lessonStr)
}

func parseState(userID int64, state *fsm.State) *subscription.RequestSubscription {
	labType, _ := state.Data["lab_type"].(polling.LabType)
	labNumber, _ := state.Data["lab_number"].(int)
	labAuditorium, _ := state.Data["lab_auditorium"].(*int)
	labDomain, _ := state.Data["lab_domain"].(*polling.LabDomain)
	labWeekday, _ := state.Data["lab_weekday"].(*int)
	labLessons, _ := state.Data["lab_lessons"].([]int)

	return &subscription.RequestSubscription{
		UserID:        int(userID),
		Type:          labType,
		LabNumber:     labNumber,
		LabAuditorium: labAuditorium,
		LabDomain:     labDomain,
		Weekday:       labWeekday,
		Lessons:       labLessons,
	}
}

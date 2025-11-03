package cmd

import (
	"context"
	"strconv"
	"strings"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *telegramBot) handleSubscriptionCreation(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	newData := &fsm.SubscriptionCreationFlowData{
		UserID: int(userID),
	}

	b.TryTransition(ctx, api, userID, fsm.StepAwaitingLabType, newData)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        askLabTypeMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: selectLabTypeKbd(),
	})
}

func (b *telegramBot) handleLabType(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	userID := update.CallbackQuery.From.ID
	labType := extractLabType(update)

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   genericServiceErrorMsg(),
		})
	}
	newData.LabType = labType

	b.TryTransition(ctx, api, userID, fsm.StepAwaitingLabNumber, newData)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      askLabNumberMsg(),
		ParseMode: models.ParseModeHTML,
	})

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
}

func (b *telegramBot) handleLabNumber(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	userID := update.Message.From.ID
	labNumberStr := update.Message.Text

	labNumber, cause := validateLabNumber(labNumberStr)
	if cause != "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   validationErrorMsg(cause),
		})
		return
	}

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   genericServiceErrorMsg(),
		})
	}
	newData.LabNumber = labNumber

	switch newData.LabType {
	case polling.LabTypePerformance:
		b.TryTransition(ctx, api, userID, fsm.StepAwaitingLabAuditorium, newData)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      askLabAuditoriumMsg(),
			ParseMode: models.ParseModeHTML,
		})
	case polling.LabTypeDefence:
		b.TryTransition(ctx, api, userID, fsm.StepAwaitingLabDomain, newData)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      userID,
			Text:        askLabDomainMsg(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: selectLabDomainKbd(),
		})
	}
}

func (b *telegramBot) handleLabAuditorium(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	userID := update.Message.From.ID
	labAuditoriumStr := update.Message.Text

	labAuditorium, cause := validateLabAuditorium(labAuditoriumStr)
	if cause != "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      validationErrorMsg(cause),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   genericServiceErrorMsg(),
		})
	}
	newData.LabAuditorium = &labAuditorium

	b.TryTransition(ctx, api, userID, fsm.StepAwaitingLabWeekday, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        askLabWeekdayMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: selectLabWeekdayKbd(),
	})
}

func (b *telegramBot) handleLabDomain(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	userID := update.CallbackQuery.From.ID
	labDomain := extractLabDomain(update)

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   genericServiceErrorMsg(),
		})
	}
	newData.LabDomain = labDomain

	b.TryTransition(ctx, api, userID, fsm.StepAwaitingLabWeekday, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        askLabWeekdayMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: selectLabWeekdayKbd(),
	})
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
}

func (b *telegramBot) handleWeekday(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	userID := update.CallbackQuery.From.ID
	weekday := extractWeekday(update)

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   genericServiceErrorMsg(),
		})
	}
	newData.Weekday = weekday

	b.TryTransition(ctx, api, userID, fsm.StepAwaitingLabLessons, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        askLabLessonsMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: selectLessonKbd(defaultLessons),
	})
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
}

func (b *telegramBot) handleLessons(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	userID := update.CallbackQuery.From.ID
	lesson := extractLesson(update)

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   genericServiceErrorMsg(),
		})
	}

	if lesson == nil {
		sub := parseFlowData(newData)
		b.TryTransition(ctx, api, userID, fsm.StepAwaitingSubCreationConfirmation, nil)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      userID,
			Text:        askSubCreationConfirmationMsg(sub),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: askLabConfirmationKbd(),
		})
		return
	}

	newData.Lessons = append(newData.Lessons, *lesson)

	existingLessonsMap := make(map[int]bool)
	for _, lesson := range newData.Lessons {
		existingLessonsMap[lesson] = true
	}

	kbdLessons := make([]Lesson, 0, len(defaultLessons))
	for i, lesson := range defaultLessons {
		lessonNum := i + 1
		if !existingLessonsMap[lessonNum] {
			kbdLessons = append(kbdLessons, lesson)
		}
	}

	b.TryTransition(ctx, api, userID, fsm.StepAwaitingLabLessons, newData)
	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      userID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: selectLessonKbd(kbdLessons),
	})
}

func (b *telegramBot) handleSubCreationConfirmation(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	userID := update.CallbackQuery.From.ID
	confirmed := extractConfirmed(update)

	newData, ok := data.(*fsm.SubscriptionCreationFlowData)
	if !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   genericServiceErrorMsg(),
		})
	}

	if confirmed {
		sub := parseFlowData(newData)
		b.TryTransition(ctx, api, userID, fsm.StepIdle, nil)

		err := b.subscriptionService.Subscribe(ctx, *sub)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    userID,
				Text:      subCreationErrorMsg(err),
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      subCreationSuccessMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	b.TryTransition(ctx, api, userID, fsm.StepIdle, nil)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      subCreationCancelledMsg(),
		ParseMode: models.ParseModeHTML,
	})
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

func extractConfirmed(update *models.Update) bool {
	confirmedStr := update.CallbackQuery.Data
	confirmedStr = strings.TrimPrefix(confirmedStr, "sub_creation:confirm:")
	return confirmedStr == "create"
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

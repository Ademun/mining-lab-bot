package cmd

import (
	"context"
	"strconv"
	"strings"

	"github.com/Ademun/mining-lab-bot/pkg/metrics"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *telegramBot) messageHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)

	state, exists := b.stateManager.get(userID)
	if !exists {
		b.defaultHandler(ctx, api, update)
		return
	}

	switch state.Step {
	case stepAwaitingLabNumber:
		b.awaitingLabNumberHandler(ctx, api, chatID, userID, text, state)
	case stepAwaitingLabAuditorium:
		b.awaitingLabAuditoriumHandler(ctx, api, chatID, userID, text, state)
	case stepAwaitingDaytime:
		b.handleAwaitingDaytime(ctx, api, chatID, userID, text, state)
	default:
		b.defaultHandler(ctx, api, update)
	}
}

func (b *telegramBot) defaultHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      startMessage(),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *telegramBot) helpHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	b.stateManager.clear(userID)
	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      helpMessage(),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *telegramBot) subscribeHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	b.stateManager.clear(userID)
	state := &userState{
		Step: stepAwaitingLabNumber,
		Data: subscriptionData{},
	}
	b.stateManager.set(userID, state)

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      subAskLabNumberMessage(),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *telegramBot) unsubscribeHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	b.stateManager.clear(userID)

	subs, err := b.subscriptionService.FindSubscriptionsByUserID(ctx, int(userID))
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      subsFetchingErrorMessage(err),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if len(subs) == 0 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      unsubEmptyListMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        unsubSelectMessage(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: createUnsubKeyboard(subs),
	})
}

func (b *telegramBot) listHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	b.stateManager.clear(userID)

	subs, err := b.subscriptionService.FindSubscriptionsByUserID(ctx, int(userID))
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      subsFetchingErrorMessage(err),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if len(subs) == 0 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      listEmptySubsMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      listSubsSuccessMessage(subs),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *telegramBot) statsHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	b.stateManager.clear(userID)

	if int(userID) != b.options.AdminID {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      permissionDeniedErrorMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	snapshot := metrics.Global().Snapshot()

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      statsSuccessMessage(&snapshot),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *telegramBot) awaitingLabNumberHandler(ctx context.Context, api *bot.Bot, chatID, userID int64, text string, state *userState) {
	labNumber, err := strconv.Atoi(text)
	if err != nil || labNumber < 1 || labNumber > 999 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      subLabNumberValidationErrorMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	state.Data.LabNumber = labNumber
	state.Step = stepAwaitingLabAuditorium
	b.stateManager.set(userID, state)

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      subAskAuditoriumMessage(),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *telegramBot) awaitingLabAuditoriumHandler(ctx context.Context, api *bot.Bot, chatID, userID int64, text string, state *userState) {
	auditorium, err := strconv.Atoi(text)
	if err != nil || auditorium < 1 || auditorium > 999 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      subAuditoriumNumberValidationErrorMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	state.Data.LabAuditorium = auditorium
	state.Step = stepAwaitingWeekday
	b.stateManager.set(userID, state)

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        subAskWeekdayMessage(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: createWeekdayKeyboard(),
	})
}

func (b *telegramBot) handleAwaitingDaytime(ctx context.Context, api *bot.Bot, chatID, userID int64, text string, state *userState) {
	_, err := parseTime(text)
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      subTimeValidationErrorMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	state.Data.Daytime = text
	state.Step = stepConfirming
	b.stateManager.set(userID, state)

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        subAskTeacherMessage(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: createSkipKeyboard("teacher"),
	})
}

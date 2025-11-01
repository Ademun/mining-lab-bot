package cmd

import (
	"context"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// /help command, or any unmatched message
func handleDefault(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.Chat.ID

	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      helpMsg(),
		ParseMode: models.ParseModeHTML,
	}); err != nil {
		slog.Error("Failed to send help message",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

// /start command
func (b *telegramBot) handleStart(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.Chat.ID

	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      startMsg(),
		ParseMode: models.ParseModeHTML,
	}); err != nil {
		slog.Error("Failed to send start message",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

// /sub command
func (b *telegramBot) handleCreatingSubscription(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	newData := map[string]interface{}{
		"user_id": userID,
	}

	slog.Info("Handling subscription creation",
		"user_id", userID,
		"service", logger.TelegramBot)

	b.tryTransition(ctx, api, userID, fsm.StepAwaitingLabType, newData)

	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        askLabTypeMsg(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: selectLabTypeKbd(),
	}); err != nil {
		slog.Error("Failed to send lab type request",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

// /unsub and /list commands
// This function is shared between two commands to provide better UX
func (b *telegramBot) handleListingSubscriptions(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	slog.Info("Handling listing / deleting subscriptions",
		"user_id", userID,
		"service", logger.TelegramBot)

	userSubs, err := b.subscriptionService.FindSubscriptionsByUserID(ctx, int(userID))
	if err != nil {
		if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      genericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		}); err != nil {
			slog.Error("Failed to send generic error message",
				"error", err,
				"user_id", userID,
				"service", logger.TelegramBot)
		}
		return
	}

	if len(userSubs) == 0 {
		if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      emptySubsListMsg(),
			ParseMode: models.ParseModeHTML,
		}); err != nil {
			slog.Error("Failed to send empty subs list message",
				"error", err,
				"user_id", userID,
				"service", logger.TelegramBot)
		}
		return
	}

	newData := map[string]interface{}{
		"user_subs": userSubs,
	}

	b.tryTransition(ctx, api, userID, fsm.StepAwaitingListingSubsAction, newData)

	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        viewSubResponseMsg(&userSubs[0]),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: listSubsKbd(userSubs[0].UUID, 0, len(userSubs)),
	}); err != nil {
		slog.Error("Failed to send subs list message",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

// /stats command, admin only
func (b *telegramBot) handleStats(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.Chat.ID
	adminID := b.options.AdminID

	if int(userID) != adminID {
		if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      permissionDeniedErrorMsg(),
			ParseMode: models.ParseModeHTML,
		}); err != nil {
			slog.Error("Failed to send permission denied message",
				"error", err,
				"user_id", userID,
				"service", logger.TelegramBot)
		}
		return
	}

	snapshot := metrics.Global()

	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      statsMsg(snapshot),
		ParseMode: models.ParseModeHTML,
	}); err != nil {
		slog.Error("Failed to send permission denied message",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
	}
}

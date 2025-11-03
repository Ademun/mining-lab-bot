package cmd

import (
	"context"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/pkg/logger"
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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      startMsg(),
		ParseMode: models.ParseModeHTML,
	})
}

package cmd

import (
	"context"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/cmd/internal/presentation"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// /help command, or any unmatched message
func handleDefault(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.Chat.ID

	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      presentation.HelpCmdMsg(),
		ParseMode: models.ParseModeHTML,
	}); err != nil {
		slog.Error("Failed to send help message",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   presentation.GenericServiceErrorMsg(),
		})
	}
}

// /start command
func (b *telegramBot) handleStart(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.Chat.ID

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      presentation.StartCmdMsg(),
		ParseMode: models.ParseModeHTML,
	})
}

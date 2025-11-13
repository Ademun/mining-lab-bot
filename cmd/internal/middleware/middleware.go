package middleware

import (
	"context"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func CommandLoggingMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil && strings.HasPrefix(update.Message.Text, "/") {
			command := strings.Fields(update.Message.Text)[0]
			recordCommand(command)
			slog.Info("Received command",
				"command", command,
				"chat_id", update.Message.Chat.ID,
				"user_id", update.Message.From.ID)
		}
		next(ctx, b, update)
	}
}

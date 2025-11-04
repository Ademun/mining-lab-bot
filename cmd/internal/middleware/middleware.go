package middleware

import (
	"context"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func TypingMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var chatID int64
		if update.Message != nil {
			chatID = update.Message.Chat.ID
		} else if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Message.Chat.ID
		} else {
			return
		}
		b.SendChatAction(ctx, &bot.SendChatActionParams{
			ChatID: chatID,
			Action: models.ChatActionTyping,
		})
		next(ctx, b, update)
	}
}

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

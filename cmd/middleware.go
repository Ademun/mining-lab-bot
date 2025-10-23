package cmd

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func typingMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var chatID int64
		if update.Message != nil {
			chatID = update.Message.Chat.ID
		} else {
			chatID = update.CallbackQuery.Message.Message.Chat.ID
		}
		b.SendChatAction(ctx, &bot.SendChatActionParams{
			ChatID: chatID,
			Action: models.ChatActionTyping,
		})
		next(ctx, b, update)
	}
}

package cmd

import (
	"context"

	"github.com/Ademun/mining-lab-bot/pkg/metrics"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// /stats command, admin only
func (b *telegramBot) handleStats(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.Chat.ID
	adminID := b.options.AdminID

	if int(userID) != adminID {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      permissionDeniedErrorMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	snapshot := metrics.Global()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      statsMsg(snapshot),
		ParseMode: models.ParseModeHTML,
	})
}

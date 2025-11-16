package cmd

import (
	"context"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/cmd/internal/presentation"
	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *telegramBot) SendNotification(ctx context.Context, notif notification.Notification) {
	userID := notif.UserID

	b.TryTransition(ctx, int64(userID), fsm.StepAwaitingFeedbackReaction, &fsm.IdleData{})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.NotifyMsg(&notif),
		ReplyMarkup: presentation.LinkKbd(notif.Slot.URL),
		ParseMode:   models.ParseModeHTML,
	})
}

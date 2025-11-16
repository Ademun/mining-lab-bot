package cmd

import (
	"context"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/cmd/internal/middleware"
	"github.com/Ademun/mining-lab-bot/cmd/internal/presentation"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *telegramBot) handleFeedbackMsg(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	b.TryTransition(ctx, userID, fsm.StepAwaitingFeedbackMsg, &fsm.IdleData{})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      presentation.FeedbackCmdMsg(),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *telegramBot) handleFeedbackMsgText(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID

	b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    b.options.AdminID,
		Text:      presentation.FeedbackRedirectMsg(userID, update.Message.Text),
		ParseMode: models.ParseModeHTML,
	})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      presentation.FeedbackReplyMsg(),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *telegramBot) handleFeedbackReaction(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if update.MessageReaction == nil {
		return
	}
	userID := update.MessageReaction.MessageID

	b.TryTransition(ctx, int64(userID), fsm.StepIdle, &fsm.IdleData{})
	for _, reaction := range update.MessageReaction.NewReaction {
		if reaction.ReactionTypeEmoji == nil {
			continue
		}
		if reaction.ReactionTypeEmoji.Emoji == "👍" {
			middleware.RecordNotification()
		}
	}
}

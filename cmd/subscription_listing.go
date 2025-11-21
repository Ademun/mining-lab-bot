package cmd

import (
	"context"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/cmd/internal/presentation"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// /unsub and /list commands
// This function is shared between two commands to provide better UX
func (b *telegramBot) handleListingSubscriptions(ctx context.Context, api *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID

	userSubs, err := b.subscriptionService.FindSubscriptionsByUserID(ctx, int(userID))
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.GenericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if len(userSubs) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.EmptySubListMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	newData := &fsm.SubscriptionListingFlowData{
		UserSubs: userSubs,
	}

	b.TryTransition(ctx, userID, fsm.StepAwaitingListingSubsAction, newData)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        presentation.SubViewMsg(&userSubs[0]),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: presentation.ListSubsKbd(userSubs[0].UUID, 0, len(userSubs)),
	})
}

func (b *telegramBot) handleListingSubsAction(ctx context.Context, api *bot.Bot, update *models.Update, data fsm.StateData) {
	if update.CallbackQuery == nil {
		return
	}
	userID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.Message.ID

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	newData, ok := data.(*fsm.SubscriptionListingFlowData)
	if !ok {
		slog.Error("Critical error: unable to assert flow data",
			"data", data,
			"service", logger.TelegramBot)
		b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.GenericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	newIndex, subUUID := extractListingData(update)
	if newIndex != nil && *newIndex >= 0 && *newIndex < len(newData.UserSubs) {
		b.api.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userID,
			MessageID: messageID,
			Text:      presentation.SubViewMsg(&newData.UserSubs[*newIndex]),
			ParseMode: models.ParseModeHTML,
		})
		b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      userID,
			MessageID:   messageID,
			ReplyMarkup: presentation.ListSubsKbd(newData.UserSubs[*newIndex].UUID, *newIndex, len(newData.UserSubs)),
		})
		return
	}

	if subUUID != nil {
		if err := b.subscriptionService.Unsubscribe(ctx, *subUUID); err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    userID,
				Text:      presentation.GenericServiceErrorMsg(),
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		newIdx := 0
		for i, sub := range newData.UserSubs {
			if sub.UUID.String() == subUUID.String() {
				newIdx = i - 1
				if newIdx < 0 {
					newIdx = 0
				}
				newData.UserSubs = append(newData.UserSubs[:i], newData.UserSubs[i+1:]...)
				break
			}
		}

		if len(newData.UserSubs) == 0 {
			b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    userID,
				MessageID: messageID,
				Text:      presentation.EmptySubListMsg(),
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userID,
			MessageID: messageID,
			Text:      presentation.SubViewMsg(&newData.UserSubs[newIdx]),
			ParseMode: models.ParseModeHTML,
		})
		b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      userID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			ReplyMarkup: presentation.ListSubsKbd(newData.UserSubs[newIdx].UUID, newIdx, len(newData.UserSubs)),
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.UnsubSuccessMsg(),
			ParseMode: models.ParseModeHTML,
		})
	}
	b.TryTransition(ctx, userID, fsm.StepAwaitingListingSubsAction, newData)
}

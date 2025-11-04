package cmd

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/cmd/internal/presentation"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

// /unsub and /list commands
// This function is shared between two commands to provide better UX
func (b *telegramBot) handleListingSubscriptions(ctx context.Context, api *bot.Bot, update *models.Update) {
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
	userID := update.CallbackQuery.From.ID

	newData, ok := data.(*fsm.SubscriptionListingFlowData)
	if !ok {
		slog.Error("Critical error: unable to assert flow data",
			"data", data,
			"service", logger.TelegramBot)
		b.TryTransition(ctx, userID, fsm.StepIdle, &fsm.IdleData{})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   presentation.GenericServiceErrorMsg(),
		})
	}

	newIndex, subUUID := extractListingData(update)
	if newIndex != nil && *newIndex > 0 && *newIndex < len(newData.UserSubs) {
		b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      userID,
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
				newData.UserSubs = append(newData.UserSubs[:i], newData.UserSubs[i+1:]...)
				break
			}
		}

		b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
			ChatID:      userID,
			ReplyMarkup: presentation.ListSubsKbd(newData.UserSubs[newIdx].UUID, newIdx, len(newData.UserSubs)),
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.UnsubSuccessMsg(),
			ParseMode: models.ParseModeHTML,
		})
	}
	b.TryTransition(ctx, userID, fsm.StepAwaitingListingSubsAction, newData)
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
}

// extractListingData returns new sub index if the selected action was "move:idx", and sub uuid if it was "delete"
func extractListingData(update *models.Update) (*int, *uuid.UUID) {
	dataFields := strings.Split(update.CallbackQuery.Data, ":")[1:]
	switch dataFields[0] {
	case "move":
		newIndex, err := strconv.Atoi(dataFields[1])
		if err != nil {
			slog.Error("Failed to parse new sub index",
				"index", dataFields[1],
				"error", err,
				"service", logger.TelegramBot)
		}
		return &newIndex, nil
	case "delete":
		subUUID, err := uuid.Parse(dataFields[1])
		if err != nil {
			slog.Error("Failed to parse sub uuid",
				"uuid", dataFields[1],
				"error", err,
				"service", logger.TelegramBot)
		}
		return nil, &subUUID
	}
	return nil, nil
}

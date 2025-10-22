package cmd

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

func (b *telegramBot) callbackRouter(ctx context.Context, api *bot.Bot, update *models.Update) {
	callbackData := update.CallbackQuery.Data
	switch {
	case strings.HasPrefix(callbackData, "weekday:"):
		b.callbackWeekdayHandler(ctx, api, update)
	case strings.HasPrefix(callbackData, "skip:"):
		b.callbackSkipHandler(ctx, api, update)
	case strings.HasPrefix(callbackData, "confirm:"):
		b.callbackConfirmSubHandler(ctx, api, update)
	case strings.HasPrefix(callbackData, "unsub:"):
		b.callbackUnsubHandler(ctx, api, update)
	}
}

func (b *telegramBot) callbackWeekdayHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	state, exists := b.stateManager.get(userID)
	if !exists || state.Step != stepAwaitingWeekday {
		api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Состояние устарело",
		})
		return
	}

	weekdayString := strings.TrimPrefix(update.CallbackQuery.Data, "weekday:")
	weekdayInt, _ := strconv.Atoi(weekdayString)
	weekday := time.Weekday(weekdayInt)

	state.Data.Weekday = &weekday
	state.Step = stepAwaitingDaytime
	b.stateManager.set(userID, state)

	api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      subAskTimeMessage(),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *telegramBot) callbackSkipHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	state, exists := b.stateManager.get(userID)
	if !exists {
		api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Состояние устарело",
		})
		return
	}

	api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	field := strings.TrimPrefix(update.CallbackQuery.Data, "skip:")
	switch field {
	case "weekday":
		if state.Step != stepAwaitingWeekday {
			return
		}
		state.Data.Weekday = nil
		state.Data.Daytime = ""
		state.Step = stepConfirming
		b.stateManager.set(userID, state)
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatID,
			Text:        subConfirmationMessage(&state.Data),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: createConfirmationKeyboard(),
		})
	}
}

func (b *telegramBot) callbackConfirmSubHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Message.Chat.ID

	action := strings.TrimPrefix(update.CallbackQuery.Data, "confirm:")
	api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	if action == "cancel" {
		b.stateManager.clear(userID)
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      subCancelledMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	state, exists := b.stateManager.get(userID)
	if !exists || state.Step != stepConfirming {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      "Состояние устарело",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	sub := model.Subscription{
		UUID:          uuid.New().String(),
		UserID:        int(userID),
		ChatID:        int(chatID),
		LabNumber:     state.Data.LabNumber,
		LabAuditorium: state.Data.LabAuditorium,
		Weekday:       state.Data.Weekday,
		DayTime:       &state.Data.Daytime,
	}

	if state.Data.Daytime != "" {
		_, err := parseTime(state.Data.Daytime)
		if err != nil {
			api.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      subCreationErrorMessage(err),
				ParseMode: models.ParseModeHTML,
			})
			b.stateManager.clear(userID)
			return
		}
	}

	if err := b.subscriptionService.Subscribe(ctx, sub); err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      subCreationErrorMessage(err),
			ParseMode: models.ParseModeHTML,
		})
		b.stateManager.clear(userID)
		return
	}

	b.stateManager.clear(userID)

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      subCreationSuccessMessage(sub.LabNumber, sub.LabAuditorium),
		ParseMode: models.ParseModeHTML,
	})

	b.notifService.NotifyNewSubscription(ctx, sub)
}

func (b *telegramBot) callbackUnsubHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	messageID := update.CallbackQuery.Message.Message.ID
	data := update.CallbackQuery.Data

	api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	switch {
	case strings.HasPrefix(data, "unsub:delete:"):
		subUUID := strings.TrimPrefix(data, "unsub:delete:")
		subs, err := b.subscriptionService.FindSubscriptionsByUserID(ctx, int(userID))
		if err != nil {
			api.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      subsFetchingErrorMessage(err),
				ParseMode: models.ParseModeHTML,
			})
		}

		var targetSub *model.Subscription
		for _, sub := range subs {
			if sub.UUID == subUUID {
				targetSub = &sub
				break
			}
		}

		if targetSub == nil {
			return
		}

		if err := b.subscriptionService.Unsubscribe(ctx, subUUID); err != nil {
			api.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      unsubErrorMessage(err),
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		updatedSubs, err := b.subscriptionService.FindSubscriptionsByUserID(ctx, int(userID))
		if err != nil {
			api.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      subsFetchingErrorMessage(err),
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		if len(updatedSubs) == 0 {
			api.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    chatID,
				MessageID: messageID,
				Text:      unsubEmptyListMessage(),
				ParseMode: models.ParseModeHTML,
			})
			return
		}
		api.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatID,
			MessageID:   messageID,
			Text:        unsubSelectMessage(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: createUnsubKeyboard(updatedSubs),
		})

		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text: unsubSuccessMessage(targetSub.LabNumber,
				targetSub.LabAuditorium),
			ParseMode: models.ParseModeHTML,
		})
	case data == "unsub:all":
		api.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatID,
			MessageID:   messageID,
			Text:        unsubConfirmDeleteAllMessage(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: createDeleteAllConfirmKeyboard(),
		})
	case data == "unsub:all:confirm":
		subs, err := b.subscriptionService.FindSubscriptionsByUserID(ctx, int(userID))
		if err != nil {
			api.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      subsFetchingErrorMessage(err),
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		count := len(subs)
		for _, sub := range subs {
			if err := b.subscriptionService.Unsubscribe(ctx, sub.UUID); err != nil {
				log.Printf("error unsubscribing: %v", err)
			}
		}

		api.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: messageID,
			Text:      unsubDeleteAllSuccessMessage(count),
			ParseMode: models.ParseModeHTML,
		})
	case data == "unsub:all:cancel":
		subs, err := b.subscriptionService.FindSubscriptionsByUserID(ctx, int(userID))
		if err != nil {
			api.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      subsFetchingErrorMessage(err),
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		api.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatID,
			MessageID:   messageID,
			Text:        unsubSelectMessage(),
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: createUnsubKeyboard(subs),
		})

	case strings.HasPrefix(data, "unsub:view:"):
		return
	}
}

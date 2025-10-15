package cmd

import (
	"context"
	"strconv"
	"strings"

	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

func (b *Bot) helpHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      helpMessage(),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *Bot) subscribeHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	args := strings.Split(update.Message.Text, " ")[1:]
	if len(args) != 2 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      subInvalidArgumentsMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	var labNumber, labAuditorium int

	num, err := strconv.Atoi(args[0])
	if err != nil || num < 1 || num > 999 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      subLabNumberValidationErrorMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	labNumber = num

	num, err = strconv.Atoi(args[1])
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      subAuditoriumNumberValidationErrorMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	labAuditorium = num

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	sub := model.Subscription{
		UUID:          uuid.New().String(),
		UserID:        int(userID),
		ChatID:        int(chatID),
		LabNumber:     labNumber,
		LabAuditorium: labAuditorium,
	}

	if err := b.subscriptionService.Subscribe(ctx, sub); err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      subCreationErrorMessage(err),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      subCreationSuccessMessage(labNumber, labAuditorium),
		ParseMode: models.ParseModeHTML,
	})

	b.notificationService.CheckCurrentSlots(ctx, sub)
}

func (b *Bot) unsubscribeHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	args := strings.Split(update.Message.Text, " ")[1:]
	if len(args) != 1 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      unsubInvalidArgumentsMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	subIdx, err := strconv.Atoi(args[0])
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      unsubInvalidSubNumberMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	chatID := update.Message.Chat.ID

	subs, err := b.subscriptionService.FindSubscriptionsByChatID(ctx, int(chatID))
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      subsFetchingErrorMessage(err),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if subIdx > len(subs) || subIdx < 1 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      unsubSubNumberValidationErrorMessage(len(subs)),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	targetSub := subs[subIdx-1]
	if err := b.subscriptionService.Unsubscribe(ctx, targetSub.UUID); err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      unsubErrorMessage(err),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      unsubSuccessMessage(targetSub.LabNumber, targetSub.LabAuditorium),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *Bot) listHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	subs, err := b.subscriptionService.FindSubscriptionsByChatID(ctx, int(chatID))
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      subsFetchingErrorMessage(err),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if len(subs) == 0 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      listEmptySubsMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      listSubsSuccessMessage(subs),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *Bot) statsHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	if int(update.Message.From.ID) != b.options.AdminID {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      permissionDeniedErrorMessage(),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	snapshot := metrics.Global().Snapshot()

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      statsSuccessMessage(&snapshot),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *Bot) notifyHandler(ctx context.Context, notifEvent event.NewNotificationEvent) {
	targetUser := notifEvent.Notification.ChatID
	slot := notifEvent.Notification.Slot

	b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    targetUser,
		Text:      notifySuccessMessage(&slot),
		ParseMode: models.ParseModeHTML,
	})
}

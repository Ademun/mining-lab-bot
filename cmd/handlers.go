package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

func (b *Bot) helpHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "<b>📖 Справка\n\n\n</b>" +
			"<b>📝 Подписка:\n\n</b>" +
			"<b>/sub &lt;номер лабы&gt; &lt;номер аудитории&gt;\n\n\n</b>" +
			"<b>⚙️ Управление:\n\n</b>" +
			"<b>/unsub &lt;номер подписки в списке&gt; - отписаться\n\n\n</b>" +
			"<b>/list - посмотреть подписки\n\n\n</b>",
		ParseMode: models.ParseModeHTML,
	})
}

func (b *Bot) subscribeHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	args := strings.Split(update.Message.Text, " ")[1:]
	if len(args) != 2 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Некорректные аргументы.\n\nИспользование: /sub &lt;номер лабы&gt; &lt;номер аудитории&gt;</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	var labNumber, labAuditorium int

	num, err := strconv.Atoi(args[0])
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Номер лабы должен быть числом</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	labNumber = num

	num, err = strconv.Atoi(args[1])
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Номер Аудитории должен быть числом</b>",
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
			Text:      fmt.Sprintf("<b>❌ Произошла ошибка при создании подписки:\n\n%s</b>", err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf(
			"<b>✅ Подписка создана!\n\n</b>"+
				"<b>📚 Лаба №%d\n\n</b>"+
				"<b>🚪 Аудитория №%d\n\n</b>"+
				"<b>Вы получите уведомление, когда появится нужная запись</b>",
			labNumber, labAuditorium,
		),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *Bot) unsubscribeHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	args := strings.Split(update.Message.Text, " ")[1:]
	if len(args) != 1 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Некорректные аргументы.\n\nИспользование: /unsub &lt;номер подпписки в списке&gt;\nЧтобы просмотреть список используйте команду /list</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	subIdx, err := strconv.Atoi(args[0])
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Номер подписки должен быть числом</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	chatID := update.Message.Chat.ID

	subs, err := b.subscriptionService.FindSubscriptionsByChatID(ctx, int(chatID))
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("<b>❌ Произошла ошибка при получении списка подписок:\n\n %s</b>", err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if subIdx > len(subs) || subIdx < 1 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("<b>❌ Номер подписки должен быть в диапазоне от 1 до числа ваших подписок - %d</b>", len(subs)),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	targetSub := subs[subIdx-1]
	if err := b.subscriptionService.Unsubscribe(ctx, targetSub.UUID); err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("<b>❌ Произошла ошибка при отписке:\n\n%s</b>", err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf(
			"✅ Вы больше не подписаны на лабу №%d в ауд. №%d",
			targetSub.LabNumber, targetSub.LabAuditorium,
		),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *Bot) listHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	subs, err := b.subscriptionService.FindSubscriptionsByChatID(ctx, int(chatID))
	if err != nil {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("<b>❌ Произошла ошибка при получении списка подписок:\n\n %s</b>", err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if len(subs) == 0 {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "🔍 У вас нет подписок на лабы.\n\nИспользуйте команду /sub &lt;номер лабы&gt; &lt;номер аудитории&gt;",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	entries := strings.Builder{}
	for idx, sub := range subs {
		entries.WriteString(fmt.Sprintf("<b>%d. Лаба №%d, ауд. №%d\n\n</b>", idx+1, sub.LabNumber, sub.LabAuditorium))
	}

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "<b>📋 Ваши подписки:\n\n</b>" + entries.String() + "<b>Для отписки используйте /unsub &lt;номер подписки в списке&gt;</b>",
		ParseMode: models.ParseModeHTML,
	})
}

func (b *Bot) notifyHandler(ctx context.Context, notifEvent event.NewNotificationEvent) {
	targetUser := notifEvent.Notification.ChatID
	labName, labNumber, labAuditorium, labDateTime := notifEvent.Notification.Slot.LabName, notifEvent.Notification.Slot.LabNumber, notifEvent.Notification.Slot.LabAuditorium, notifEvent.Notification.Slot.DateTime

	b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: targetUser,
		Text: fmt.Sprintf("<b>🔥 Появилась запись!\n\n\n</b>"+
			"<b>📚 Лаба №%d. %s\n\n</b>"+
			"<b>🚪 Аудитория №%d\n\n</b>"+
			"<b>🗓️ Когда: %s</b>",
			labNumber, labName, labAuditorium, formatDateTime(labDateTime)),
		ParseMode: models.ParseModeHTML,
	})
}

package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
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

	b.notificationService.CheckCurrentSlots(ctx, sub)
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

func (b *Bot) statsHandler(ctx context.Context, api *bot.Bot, update *models.Update) {
	if int(update.Message.From.ID) != b.options.AdminID {
		api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("<b>❌ Доступ запрещён. Команда доступна только для разработчика</b>"),
			ParseMode: models.ParseModeHTML,
		})
	}

	snapshot := metrics.Global().Snapshot()
	uptime := time.Since(snapshot.StartTime)

	statsText := strings.Builder{}
	statsText.WriteString("<b>📊 Статистика бота\n\n</b>")

	statsText.WriteString("<b>🕐 Общее время работы:</b> ")
	statsText.WriteString(formatDuration(uptime))
	statsText.WriteString("\n\n")

	statsText.WriteString("<b>🔍 Опросы:\n</b>")
	statsText.WriteString(fmt.Sprintf("  Всего опросов: <b>%d</b>\n", snapshot.PollingMetrics.TotalPolls))
	statsText.WriteString(fmt.Sprintf("  Режим: <b>%s</b>\n", formatPollingMode(snapshot.PollingMetrics.Mode)))
	statsText.WriteString(fmt.Sprintf("  Ошибки парсинга: <b>%d</b>\n", snapshot.PollingMetrics.ParsingErrors))
	statsText.WriteString(fmt.Sprintf("  Ошибки получения: <b>%d</b>\n", snapshot.PollingMetrics.FetchErrors))
	statsText.WriteString(fmt.Sprintf("  Среднее время опроса: <b>%s</b>\n", snapshot.PollingMetrics.AveragePollingTime.Round(time.Millisecond)))
	statsText.WriteString(fmt.Sprintf("  Среднее количество слотов: <b>%d</b>\n\n", snapshot.PollingMetrics.AverageSlotNumber))

	statsText.WriteString("<b>🔔 Уведомления:\n</b>")
	statsText.WriteString(fmt.Sprintf("  Всего уведомлений: <b>%d</b>\n", snapshot.NotificationMetrics.TotalNotifications))
	statsText.WriteString(fmt.Sprintf("  Размер кеша: <b>%d</b>\n", snapshot.NotificationMetrics.CacheLength))
	statsText.WriteString(fmt.Sprintf("  Среднее количество уведомлений: <b>%d</b>\n\n", snapshot.NotificationMetrics.AverageNotifications))

	statsText.WriteString("<b>📝 Подписки:\n</b>")
	statsText.WriteString(fmt.Sprintf("  Активных подписок: <b>%d</b>", snapshot.SubscriptionMetrics.TotalSubscriptions))

	api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      statsText.String(),
		ParseMode: models.ParseModeHTML,
	})
}

func (b *Bot) notifyHandler(ctx context.Context, notifEvent event.NewNotificationEvent) {
	targetUser := notifEvent.Notification.ChatID
	slot := notifEvent.Notification.Slot

	b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: targetUser,
		Text: fmt.Sprintf("<b>🔥 Появилась запись!\n\n\n</b>"+
			"<b>📚 Лаба №%d. %s\n\n</b>"+
			"<b>🚪 Аудитория №%d\n\n</b>"+
			"<b>🗓️ Когда: %s\n\n</b>"+
			"<b>🔗 <a href='%s'>Ссылка на запись</a></b>",
			slot.LabNumber, slot.LabName, slot.LabAuditorium, formatDateTime(slot.DateTime), slot.URL),
		ParseMode: models.ParseModeHTML,
	})
}

package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

func (bt *Bot) helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.ID,
		Text: "<b>📖 Справка\n\n\n</b>" +
			"<b>📝 Подписка:\n\n</b>" +
			"<b>/sub &lt;номер лабы&gt; &lt;номер аудитории&gt;</b>\n\n\n" +
			"<b>⚙️ Управление:\n\n</b>" +
			"<b>/unsub &lt;номер подпписки в списке&gt; - отписаться</b>\n\n\n" +
			"<b>/list - посмотреть подписки</b>\n\n\n",
		ParseMode: models.ParseModeHTML,
	})
}

func (bt *Bot) subscribeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	args := strings.Split(update.Message.Text, " ")[1:]
	if len(args) != 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Некорректные аргументы.\n\nИспользование: /sub &lt;номер лабы&gt; &lt;номер аудитории&gt;</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	var labNumber, labAuditorium int

	if num, err := strconv.Atoi(args[0]); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Номер лабы должен быть числом</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	} else {
		labNumber = num
	}

	if num, err := strconv.Atoi(args[1]); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Номер Аудитории должен быть числом</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	} else {
		labAuditorium = num
	}

	userID := update.Message.From.ID

	sub := model.Subscription{
		UUID:          uuid.New().String(),
		UserID:        int(userID),
		LabNumber:     labNumber,
		LabAuditorium: labAuditorium,
	}

	if err := bt.subService.Subscribe(ctx, sub); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("<b>❌ Произошла ошибка при создании подписки:\n\n%s</b>", err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
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

func (bt *Bot) unsubscribeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	args := strings.Split(update.Message.Text, " ")[1:]
	if len(args) != 1 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Некорректные аргументы.\n\nИспользование: /unsub &lt;номер подпписки в списке&gt;\nЧтобы просмотреть список используйте команду /list</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	subIdx, err := strconv.Atoi(args[0])
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Номер подписки должен быть числом</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	userID := update.Message.From.ID
	subs, err := bt.subService.ListForUser(ctx, int(userID))
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("<b>❌ Произошла ошибка при получении списка подписок:\n\n %s</b>", err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if subIdx > len(subs) || subIdx < 1 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("<b>❌ Номер подписки должен быть в диапазоне от 1 до числа ваших подписок - %d</b>", len(subs)),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	targetSub := subs[subIdx-1]
	if err := bt.subService.Unsubscribe(ctx, targetSub.UUID); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("<b>❌ Произошла ошибка при отписке:\n\n%s</b>", err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf(
			"✅ Вы больше не подписаны на лабу №%d в ауд. №%d",
			targetSub.LabNumber, targetSub.LabAuditorium,
		),
		ParseMode: models.ParseModeHTML,
	})
}

func (bt *Bot) listHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	subs, err := bt.subService.ListForUser(ctx, int(userID))
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("<b>❌ Произошла ошибка при получении списка подписок:\n\n %s</b>", err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if len(subs) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "<b>📋 Ваши подписки:\n\n</b>" + entries.String() + "<b>Для отписки используйте /unsub &lt;номер подписки в списке&gt;</b>",
		ParseMode: models.ParseModeHTML,
	})
}

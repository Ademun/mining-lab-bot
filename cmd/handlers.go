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
)

func (bt *Bot) subscribeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	args := strings.Split(update.Message.Text, " ")[1:]
	if len(args) != 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Недостаточно аргументов. Использование: '/sub <номер аудитории> <номер лабы>'</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	var labAuditorium, labNumber int
	if num, err := strconv.Atoi(args[0]); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Номер Аудитории должен быть числом</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	} else {
		labAuditorium = num
	}
	if num, err := strconv.Atoi(args[1]); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<b>❌ Номер лабы должен быть числом</b>",
			ParseMode: models.ParseModeHTML,
		})
		return
	} else {
		labNumber = num
	}

	userID := update.Message.From.ID

	sub := model.Subscription{
		ID:            -1,
		UserID:        int(userID),
		LabNumber:     labNumber,
		LabAuditorium: labAuditorium,
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf(
			"<b>✅ Подписка создана!\n\n</b>"+
				"<b>🚪 Аудитория №%d\n\n</b>"+
				"<b>📚 Лаба №%d\n\n</b>"+
				"<b>Вы получите уведомление, когда появится нужная запись</b>",
			labAuditorium, labNumber,
		),
		ParseMode: models.ParseModeHTML,
	})
}

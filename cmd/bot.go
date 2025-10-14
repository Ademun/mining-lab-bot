package cmd

import (
	"context"
	"fmt"

	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	eventBus            *event.Bus
	subscriptionService subscription.SubscriptionService
	notificationService notification.NotificationService
	api                 *bot.Bot
	options             *config.TelegramConfig
}

func NewBot(eb *event.Bus, subService subscription.SubscriptionService, notifService notification.NotificationService, opts *config.TelegramConfig) (*Bot, error) {
	botOpts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	b, err := bot.New(opts.BotToken, botOpts...)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %w", err)
	}

	return &Bot{
		eventBus:            eb,
		subscriptionService: subService,
		notificationService: notifService,
		api:                 b,
		options:             opts,
	}, nil
}

func (b *Bot) Start(ctx context.Context) {
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommandStartOnly, b.helpHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "sub", bot.MatchTypeCommandStartOnly, b.subscribeHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "unsub", bot.MatchTypeCommandStartOnly, b.unsubscribeHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "list", bot.MatchTypeCommandStartOnly, b.listHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "stats", bot.MatchTypeCommandStartOnly, b.statsHandler)

	event.Subscribe(b.eventBus, b.notifyHandler)

	go b.api.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "<b>👋 Привет!\n\n\n</b>" + "<b>Я бот для записи на лабораторные работы\n\n\n</b>" +
			"<b>Буду следить за появлением доступных записей и сразу уведомлю тебя, когда появится нужная\n\n\n</b>" +
			"<b>Используй /help для просмотра доступных команд</b>",
		ParseMode: models.ParseModeHTML,
	})
}

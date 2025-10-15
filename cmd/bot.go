package cmd

import (
	"context"
	"fmt"

	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot interface {
	Start(ctx context.Context)
	SendNotification(ctx context.Context, notif model.Notification)
}

type telegramBot struct {
	subscriptionService subscription.Service
	notifService        notification.Service
	api                 *bot.Bot
	options             *config.TelegramConfig
}

func NewBot(subService subscription.Service, notifService notification.Service, opts *config.TelegramConfig) (Bot, error) {
	botOpts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	b, err := bot.New(opts.BotToken, botOpts...)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %w", err)
	}

	return &telegramBot{
		subscriptionService: subService,
		notifService:        notifService,
		api:                 b,
		options:             opts,
	}, nil
}

func (b *telegramBot) Start(ctx context.Context) {
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommandStartOnly, b.helpHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "sub", bot.MatchTypeCommandStartOnly, b.subscribeHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "unsub", bot.MatchTypeCommandStartOnly, b.unsubscribeHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "list", bot.MatchTypeCommandStartOnly, b.listHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "stats", bot.MatchTypeCommandStartOnly, b.statsHandler)

	go b.api.Start(ctx)
}

func (b *telegramBot) SendNotification(ctx context.Context, notif model.Notification) {
	targetUser := notif.ChatID
	slot := notif.Slot

	b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    targetUser,
		Text:      notifySuccessMessage(&slot),
		ParseMode: models.ParseModeHTML,
	})
}

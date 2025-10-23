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
	SetNotificationService(svc notification.Service)
}

type telegramBot struct {
	subscriptionService subscription.Service
	notifService        notification.Service
	api                 *bot.Bot
	stateManager        *stateManager
	options             *config.TelegramConfig
}

func NewBot(subService subscription.Service, opts *config.TelegramConfig) (Bot, error) {
	botOpts := []bot.Option{bot.WithMiddlewares(typingMiddleware)}
	b, err := bot.New(opts.BotToken, botOpts...)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %w", err)
	}

	return &telegramBot{
		subscriptionService: subService,
		api:                 b,
		stateManager:        newStateManager(),
		options:             opts,
	}, nil
}

func (b *telegramBot) Start(ctx context.Context) {
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "help",
		bot.MatchTypeCommandStartOnly, b.helpHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "sub",
		bot.MatchTypeCommandStartOnly, b.subscribeHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "unsub",
		bot.MatchTypeCommandStartOnly, b.unsubscribeHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "list",
		bot.MatchTypeCommandStartOnly, b.listHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "stats",
		bot.MatchTypeCommandStartOnly, b.statsHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains,
		b.messageHandler)
	b.api.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix,
		b.callbackRouter)
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

func (b *telegramBot) SetNotificationService(svc notification.Service) {
	b.notifService = svc
}

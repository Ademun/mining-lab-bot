package cmd

import (
	"context"
	"fmt"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/redis/go-redis/v9"
)

type Bot interface {
	Start(ctx context.Context)
	SendNotification(ctx context.Context, notif notification.Notification)
	SetNotificationService(svc notification.Service)
}

type telegramBot struct {
	subscriptionService subscription.Service
	notifService        notification.Service
	api                 *bot.Bot
	router              *fsm.Router
	options             *config.TelegramConfig
}

func NewBot(subService subscription.Service, opts *config.TelegramConfig, redis *redis.Client) (Bot, error) {
	router := fsm.NewRouter(fsm.NewFSM(redis))
	botOpts := []bot.Option{
		bot.WithMiddlewares(typingMiddleware, router.Middleware),
		bot.WithDefaultHandler(handleDefault),
	}
	b, err := bot.New(opts.BotToken, botOpts...)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %w", err)
	}

	return &telegramBot{
		subscriptionService: subService,
		api:                 b,
		router:              router,
		options:             opts,
	}, nil
}

func (b *telegramBot) Start(ctx context.Context) {
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "start",
		bot.MatchTypeCommandStartOnly, b.handleStart)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "help",
		bot.MatchTypeCommandStartOnly, handleDefault)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "sub",
		bot.MatchTypeCommandStartOnly, b.handleCreatingSubscription)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "unsub",
		bot.MatchTypeCommandStartOnly, b.handleListingSubscriptions)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "list",
		bot.MatchTypeCommandStartOnly, b.handleListingSubscriptions)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "stats",
		bot.MatchTypeCommandStartOnly, b.handleStats)
	b.router.RegisterHandler(fsm.StepAwaitingLabType, b.handleLabType)
	b.router.RegisterHandler(fsm.StepAwaitingLabNumber, b.handleLabNumber)
	b.router.RegisterHandler(fsm.StepAwaitingLabAuditorium, b.handleLabAuditorium)
	b.router.RegisterHandler(fsm.StepAwaitingLabDomain, b.handleLabDomain)
	b.router.RegisterHandler(fsm.StepAwaitingLabWeekday, b.handleLabWeekday)
	b.router.RegisterHandler(fsm.StepAwaitingLabLessons, b.handleLabLessons)
	b.router.RegisterHandler(fsm.StepAwaitingSubCreationConfirmation, b.handleSubCreationConfirmation)
	b.router.RegisterHandler(fsm.StepAwaitingListingSubsAction, b.handleListingSubsAction)
	go b.api.Start(ctx)
}

func (b *telegramBot) SendNotification(ctx context.Context, notif notification.Notification) {
	targetUser := notif.UserID

	hidePreview := true
	b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: targetUser,
		Text:   notifySuccessMessage(&notif),
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: &hidePreview,
		},
		ReplyMarkup: createLinkKeyboard(notif.Slot.URL),
		ParseMode:   models.ParseModeHTML,
	})
}

func (b *telegramBot) SetNotificationService(svc notification.Service) {
	b.notifService = svc
}

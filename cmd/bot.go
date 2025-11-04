package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/cmd/fsm"
	"github.com/Ademun/mining-lab-bot/cmd/internal/middleware"
	"github.com/Ademun/mining-lab-bot/cmd/internal/presentation"
	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/redis/go-redis/v9"
)

type Bot interface {
	Start(ctx context.Context)
	SetNotificationService(svc notification.Service)
	SendMessage(ctx context.Context, params *bot.SendMessageParams)
	SendNotification(ctx context.Context, notif notification.Notification)
	AnswerCallbackQuery(ctx context.Context, params *bot.AnswerCallbackQueryParams)
	EditMessageReplyMarkup(ctx context.Context, params *bot.EditMessageReplyMarkupParams)
	EditMessageText(ctx context.Context, params *bot.EditMessageTextParams)
	TryTransition(ctx context.Context, userID int64, newStep fsm.ConversationStep, newData fsm.StateData)
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
		bot.WithMiddlewares(middleware.CommandLoggingMiddleware, middleware.TypingMiddleware, router.Middleware),
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
		bot.MatchTypeCommandStartOnly, b.handleSubscriptionCreation)
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
	b.router.RegisterHandler(fsm.StepAwaitingLabWeekday, b.handleWeekday)
	b.router.RegisterHandler(fsm.StepAwaitingLabLessons, b.handleLessons)
	b.router.RegisterHandler(fsm.StepAwaitingSubCreationConfirmation, b.handleSubCreationConfirmation)
	b.router.RegisterHandler(fsm.StepAwaitingListingSubsAction, b.handleListingSubsAction)
	go b.api.Start(ctx)
}

func (b *telegramBot) SetNotificationService(svc notification.Service) {
	b.notifService = svc
}

func (b *telegramBot) SendMessage(ctx context.Context, params *bot.SendMessageParams) {
	if _, err := b.api.SendMessage(ctx, params); err != nil {
		slog.Error("Failed to send message",
			"error", err,
			"params", params)
		b.api.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: params.ChatID,
			Text:   presentation.GenericServiceErrorMsg(),
		})
	}
}

func (b *telegramBot) SendNotification(ctx context.Context, notif notification.Notification) {
	targetUser := notif.UserID

	hidePreview := true
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: targetUser,
		Text:   presentation.NotifyMsg(&notif),
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: &hidePreview,
		},
		ReplyMarkup: presentation.LinkKbd(notif.Slot.URL),
		ParseMode:   models.ParseModeHTML,
	})
}

func (b *telegramBot) AnswerCallbackQuery(ctx context.Context, params *bot.AnswerCallbackQueryParams) {
	if _, err := b.api.AnswerCallbackQuery(ctx, params); err != nil {
		slog.Error("Failed to answer callback query",
			"error", err,
			"params", params)
	}
}

func (b *telegramBot) EditMessageReplyMarkup(ctx context.Context, params *bot.EditMessageReplyMarkupParams) {
	if _, err := b.api.EditMessageReplyMarkup(ctx, params); err != nil {
		slog.Error("Failed to edit message reply markup",
			"error", err,
			"params", params)
	}
}

func (b *telegramBot) EditMessageText(ctx context.Context, params *bot.EditMessageTextParams) {
	if _, err := b.api.EditMessageText(ctx, params); err != nil {
		slog.Error("Failed to edit message text",
			"error", err,
			"params", params)
	}
}

func (b *telegramBot) TryTransition(ctx context.Context, userID int64, newStep fsm.ConversationStep, newData fsm.StateData) {
	if err := b.router.Transition(ctx, userID, newStep, newData); err != nil {
		slog.Error("State transition failed",
			"error", err,
			"user_id", userID,
			"new_step", newStep,
			"service", logger.TelegramBot)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      presentation.GenericServiceErrorMsg(),
			ParseMode: models.ParseModeHTML,
		})
	}
}

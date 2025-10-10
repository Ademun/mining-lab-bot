package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

type Bot struct {
	ctx context.Context
	eb  *event.Bus
	bot *bot.Bot
}

func NewBot(ctx context.Context, eb *event.Bus) (*Bot, error) {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error(fmt.Sprintf("Error loading .env file: %v", err))
		os.Exit(1)
	}
	botKey := os.Getenv("TG_BOT_KEY")
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}
	b, err := bot.New(botKey, opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating the bot: %w", err)
	}

	return &Bot{
		ctx: ctx,
		eb:  eb,
		bot: b,
	}, nil
}

func (bt *Bot) Start() {
	bt.bot.RegisterHandler(bot.HandlerTypeMessageText, "sub", bot.MatchTypeCommandStartOnly, bt.subscribeHandler)
	go bt.bot.Start(bt.ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

}

package main

import (
	"context"
	"log"
	"time"

	"github.com/Ademun/mining-lab-bot/cmd"
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/pkg/event"
)

func main() {
	eb := event.NewEventBus()
	ctx := context.Background()

	ps := polling.New(eb, nil)
	if err := ps.Start(ctx); err != nil {
		log.Fatal(err)
	}

	bot, err := cmd.NewBot(ctx, eb)
	if err != nil {
		log.Fatal(err)
	}
	bot.Start()

	time.Sleep(1 * time.Hour)
}

package main

import (
	"context"
	"log"
	"time"

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

	time.Sleep(1 * time.Hour)
}

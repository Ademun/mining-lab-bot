package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ademun/mining-lab-bot/cmd"
	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/event"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := sql.Open("sqlite3", "./dev.db")
	if err != nil {
		log.Fatal(err)
	}

	subRepo, err := subscription.NewRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	eb := event.NewEventBus()

	ps := polling.New(eb, nil)
	if err := ps.Start(ctx); err != nil {
		log.Fatal(err)
	}

	ss := subscription.New(eb, subRepo)
	if err := ss.Start(ctx); err != nil {
		log.Fatal(err)
	}

	ns := notification.New(eb, ss)
	if err := ns.Start(ctx); err != nil {
		log.Fatal(err)
	}

	bot, err := cmd.NewBot(ctx, eb, ss)
	if err != nil {
		log.Fatal(err)
	}
	bot.Start()

	<-ctx.Done()
	slog.Info("Shutting down...")
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	db.Close()
}

package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ademun/mining-lab-bot/cmd"
	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/internal/teacher"
	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger.Init(slog.LevelInfo)

	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("Fatal error", "error", err)
		return
	}

	db, err := sql.Open("sqlite3", "./dev.db")
	if err != nil {
		slog.Error("Fatal error", "error", err)
		return
	}

	subscriptionRepo, err := subscription.NewRepo(ctx, db)
	if err != nil {
		slog.Error("Fatal error", "error", err)
		return
	}

	subscriptionService := subscription.New(subscriptionRepo)
	if err := subscriptionService.Start(ctx); err != nil {
		slog.Error("Fatal error", "error", err)
	}

	bot, err := cmd.NewBot(subscriptionService, &cfg.TelegramConfig)
	if err != nil {
		slog.Error("Fatal error", "error", err)
		return
	}

	notificationService := notification.New(subscriptionService, bot)

	bot.SetNotificationService(notificationService)
	bot.Start(ctx)

	teacherRepo, err := teacher.NewRepo(ctx, db)
	if err != nil {
		slog.Error("Fatal error", "error", err)
		return
	}
	teacherService := teacher.New(teacherRepo)

	pollingService := polling.New(notificationService, teacherService, &cfg.PollingConfig)
	if err := pollingService.Start(ctx); err != nil {
		slog.Error("Fatal error", "error", err)
		return
	}

	<-ctx.Done()
	slog.Info("Shutting down...")
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := db.Close(); err != nil {
		slog.Error("Fatal error", "error", err)
		return
	}
}

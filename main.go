package main

import (
	"context"
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
	"github.com/jmoiron/sqlx"
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

	db, err := sqlx.Open("sqlite3", "./dev.db")
	if err != nil {
		slog.Error("Fatal error", "error", err)
		return
	}

	subscriptionRepo := subscription.NewRepo(db)

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

	teacherRepo := teacher.NewRepo(db)
	teacherService := teacher.New(teacherRepo, &cfg.TeacherConfig)

	pollingService := polling.New(notificationService, teacherService, &cfg.PollingConfig)
	if err := pollingService.Start(ctx); err != nil {
		slog.Error("Fatal error", "error", err)
		return
	}

	<-ctx.Done()
	slog.Info("Shutting down...")
	ctx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	pollingService.Stop(ctx)

	if err := db.Close(); err != nil {
		slog.Error("Fatal error", "error", err)
		return
	}
}

package logger

import (
	"log/slog"
	"os"
)

const (
	ServicePolling      = "polling"
	ServiceNotification = "notification"
	ServiceSubscription = "subscription"
	TelegramBot         = "bot"
)

func Init(level slog.Level) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))
}

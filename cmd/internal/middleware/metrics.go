package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	commandUsageMetrics = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "telegram_command_usage",
		Help: "Command usage",
	}, []string{"command"})

	notificationConversionMetrics = promauto.NewCounter(prometheus.CounterOpts{
		Name: "telegram_notification_conversion",
		Help: "Notification conversion",
	})
)

func recordCommand(command string) {
	commandUsageMetrics.WithLabelValues(command).Inc()
}

func RecordNotification() {
	notificationConversionMetrics.Inc()
}

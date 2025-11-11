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

	linkClicksMetrics = promauto.NewCounter(prometheus.CounterOpts{
		Name: "telegram_link_clicks",
		Help: "Link clicks",
	})
)

func recordCommand(command string) {
	commandUsageMetrics.WithLabelValues(command).Inc()
}

func RecordLinkClick() {
	linkClicksMetrics.Inc()
}

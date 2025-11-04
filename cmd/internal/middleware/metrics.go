package middleware

import "github.com/prometheus/client_golang/prometheus"

var (
	commandUsageMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "telegram_command_usage",
		Help: "Command usage",
	}, []string{"command"})
)

func recordCommand(command string) {
	commandUsageMetrics.WithLabelValues(command).Inc()
}

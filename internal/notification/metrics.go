package notification

import "github.com/prometheus/client_golang/prometheus"

var (
	notificationsSentMetics = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "notifications_sent",
		Help: "Number of notifications sent",
	})
)

func recordNotification() {
	notificationsSentMetics.Inc()
}

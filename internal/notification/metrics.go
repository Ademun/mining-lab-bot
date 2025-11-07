package notification

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	notificationsSentMetics = promauto.NewCounter(prometheus.CounterOpts{
		Name: "notifications_sent",
		Help: "Number of notifications sent",
	})
)

func recordNotification() {
	notificationsSentMetics.Inc()
}

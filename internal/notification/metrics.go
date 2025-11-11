package notification

import (
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	notificationsSentMetics = promauto.NewCounter(prometheus.CounterOpts{
		Name: "notifications_sent",
		Help: "Number of notifications sent",
	})

	uniqueSlotsMetrics = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "notifications_unique_slots",
		Help: "Slot count by type",
	}, []string{"type"})
)

func recordNotification() {
	notificationsSentMetics.Inc()
}

func recordSlot(slotType polling.LabType) {
	var enType string
	switch slotType {
	case polling.LabTypePerformance:
		enType = "performance"
	case polling.LabTypeDefence:
		enType = "defence"
	}
	uniqueSlotsMetrics.WithLabelValues(enType).Inc()
}

package polling

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "polling_http_requests_metrics",
		Help: "HTTP request duration in seconds by status code",
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.95: 0.01,
			0.99: 0.001,
		},
	}, []string{"status_code"})

	parsingErrorsCountMetrics = promauto.NewCounter(prometheus.CounterOpts{
		Name: "polling_parsing_errors_count",
		Help: "Number of parsing errors",
	})

	parsingDurationMetrics = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "polling_parsing_duration",
		Help: "Duration in seconds of parsed slots",
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.95: 0.01,
			0.99: 0.001,
		},
	})

	pollingDurationMetrics = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "polling_duration",
		Help: "Slot polling duration in seconds",
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.95: 0.01,
			0.99: 0.001,
		},
	})

	slotCountMetrics = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "polling_slot_count",
		Help: "Slot count by type",
	}, []string{"type"})
)

func recordRequest(d time.Duration, statusCode int) {
	httpRequestsMetrics.WithLabelValues(http.StatusText(statusCode)).Observe(d.Seconds())
}

func recordParsing(d time.Duration, hasError bool) {
	if hasError {
		parsingErrorsCountMetrics.Inc()
	}
	parsingDurationMetrics.Observe(d.Seconds())
}

func recordPolling(d time.Duration) {
	pollingDurationMetrics.Observe(d.Seconds())
}

func recordSlot(slotType LabType) {
	var enType string
	switch slotType {
	case LabTypePerformance:
		enType = "performance"
	case LabTypeDefence:
		enType = "defence"
	}
	slotCountMetrics.WithLabelValues(enType).Inc()
}

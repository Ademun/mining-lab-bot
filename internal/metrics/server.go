package metrics

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(addr string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	go func() {
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			slog.Error("Metrics server error", "error", err)
			return
		}
	}()
}

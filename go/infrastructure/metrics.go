package infrastructure

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MustStartMetrics() {
	log.Info("Starting metrics server at :9090/metrics")
	metricsserver := http.NewServeMux()
	metricsserver.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Handler: metricsserver,
		Addr:    ":9090",
	}
	if err := server.ListenAndServe() ; err != nil {
		panic(err)
	}
}


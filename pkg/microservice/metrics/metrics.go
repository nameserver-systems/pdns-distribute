package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartMetricsExporter is blocking.
func StartMetricsExporter(address string) error {
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.MustRegister(collectors.NewGoCollector(collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsAll)))
	prometheus.MustRegister(collectors.NewBuildInfoCollector())

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}))

	server := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

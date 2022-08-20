package metrics

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartMetricsExporter is blocking.
func StartMetricsExporter(address string) error {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Handle("/metrics", promhttp.Handler())

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

package servicediscovery

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/servicediscovery"
)

var (
	servicediscoverycalltotal = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_service_discovery_call_total", ConstLabels: map[string]string{"state": "successful"}, Help: "The total count of service discovery calls"})
	servicediscoveryerrtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_service_discovery_call_total", ConstLabels: map[string]string{"state": "failed"}, Help: "The total count of service discovery calls"})
)

func GetActiveSecondaries(ms *microservice.Microservice) ([]servicediscovery.ResolvedService, error) {
	services, resolveerr := ms.ServiceDiscovery.GetServices("pdns-secondary-syncer", "")
	if resolveerr != nil {
		servicediscoveryerrtotal.Inc()

		return nil, resolveerr
	}

	servicediscoverycalltotal.Inc()

	return services, nil
}

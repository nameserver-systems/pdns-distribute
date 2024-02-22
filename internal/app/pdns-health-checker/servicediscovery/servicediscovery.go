package servicediscovery

import (
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/messaging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	servicediscoverycalltotal = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_service_discovery_call_total", ConstLabels: map[string]string{"state": "successful"}, Help: "The total count of service discovery calls"})
	servicediscoveryerrtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_service_discovery_call_total", ConstLabels: map[string]string{"state": "failed"}, Help: "The total count of service discovery calls"})
)

func GetActiveSecondaries(ms *microservice.Microservice) ([]messaging.ResolvedService, error) {
	stream := ms.MessageBroker.GetStream()

	services, resolveerr := ms.MessageBroker.RetrieveRegisteredConsumers(stream)
	if resolveerr != nil {
		servicediscoveryerrtotal.Inc()

		return nil, resolveerr
	}

	servicediscoverycalltotal.Inc()

	return services, nil
}

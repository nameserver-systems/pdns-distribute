package messaging

import (
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-zone-provider/config"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-zone-provider/powerdns"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requesttotal = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_request_total", Help: "The total count of requests"})

func SubscribeToIncomingMessages(ms *microservice.Microservice, serviceconfig *config.ServiceConfiguration) {
	topic := serviceconfig.ReceiverTopic
	queue := "worker"

	ms.MessageBroker.SubscribeQueueAsync(topic, queue, func(msg *nats.Msg) {
		go powerdns.ProcessRequest(msg, serviceconfig)
		requesttotal.Inc()
	})
}

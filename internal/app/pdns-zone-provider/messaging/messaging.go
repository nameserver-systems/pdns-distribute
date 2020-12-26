package messaging

import (
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-zone-provider/config"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-zone-provider/powerdns"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice"
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

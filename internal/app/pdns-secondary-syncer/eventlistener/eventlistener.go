package eventlistener

import (
	"strings"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/config"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/modeljob"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/powerdns"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/primaryzoneprovider"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/worker"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	changereceivedtotal   = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_received_zone_events_total", ConstLabels: map[string]string{"event_type": "change"}, Help: "The total count of received zone events"})
	deletereceivedtotal   = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_received_zone_events_total", ConstLabels: map[string]string{"event_type": "delete"}, Help: "The total count of received zone events"})
	createreceivedtotal   = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_received_zone_events_total", ConstLabels: map[string]string{"event_type": "create"}, Help: "The total count of received zone events"})
	zonestaterequesttotal = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_zone_state_request_total", Help: "The total count of zone state requests"})
)

func StartEventListenerAndWorker(ms *microservice.Microservice) error {
	serviceconfig := config.GetConfiguration(ms)

	port, err := ms.GetServicePort()
	if err != nil {
		return err
	}

	stream, err := ms.MessageBroker.CreatePersistentMessageStore("pdns-distribute-event-store", []string{"zone.>"})
	if err != nil {
		return err
	}

	consumer, err := ms.MessageBroker.CreatePersistentMessageReceiver("secondary-syncer-event-client", ms.ID, "add", port, "secondary", stream)
	if err != nil {
		return err
	}

	ms.MessageBroker.SetStream(stream)
	ms.MessageBroker.SetConsumer(consumer)

	startZoneEventListeners(ms, serviceconfig)
	startSecondaryZoneStateEventListener(ms, serviceconfig)
	primaryzoneprovider.StartAXFRProvider(ms, serviceconfig)

	worker.StartWorker(serviceconfig)

	return nil
}

func StopWorker() {
	worker.CloseWorkerQueue()
}

func StopDNSServer() {
	if err := primaryzoneprovider.StopDNSServer(); err != nil {
		logger.ErrorErrLog(err)
	}
}

func startZoneEventListeners(ms *microservice.Microservice, conf *config.ServiceConfiguration) (cancelCtx jetstream.ConsumeContext) {
	addtopic := conf.AddEventTopic
	changetopic := conf.ChangeEventTopic
	deltopic := conf.DeleteEventTopic

	consumer := ms.MessageBroker.GetConsumer()
	cancelCtx, err := consumer.Consume(func(msg jetstream.Msg) {
		subject := msg.Subject()
		switch {
		case strings.HasPrefix(subject, addtopic):
			go worker.EnqueJob(&modeljob.PowerDNSAPIJob{
				Jobtype: modeljob.AddZone,
				Msg:     msg,
				Ms:      ms,
				Conf:    conf,
			})
			createreceivedtotal.Inc()
		case strings.HasPrefix(subject, changetopic):
			go worker.EnqueJob(&modeljob.PowerDNSAPIJob{
				Jobtype: modeljob.ChangeZone,
				Msg:     msg,
				Ms:      ms,
				Conf:    conf,
			})
			changereceivedtotal.Inc()
		case strings.HasPrefix(subject, deltopic):
			go worker.EnqueJob(&modeljob.PowerDNSAPIJob{
				Jobtype: modeljob.DeleteZone,
				Msg:     msg,
				Ms:      ms,
				Conf:    conf,
			})
			deletereceivedtotal.Inc()
		default:
			logger.DebugLog("[Zone Event Listener]: not matched on topic: " + subject)
		}
		if err := msg.Ack(); err != nil {
			logger.ErrorErrLog(err)
		}
	})
	if err != nil {
		logger.ErrorErrLog(err)
	}

	return
}

func startSecondaryZoneStateEventListener(ms *microservice.Microservice, conf *config.ServiceConfiguration) {
	receivertopic := conf.ZoneStateReceiverTopic

	zonestatehandler := func(msg *nats.Msg) {
		go powerdns.RespondSecondaryZoneState(msg, conf)
		logger.DebugLog("[Zone State Request] requested with query: " + string(msg.Data) + " from: " + msg.Reply)
		zonestaterequesttotal.Inc()
	}

	if err := ms.MessageBroker.SubscribeAsync(receivertopic, zonestatehandler); err != nil {
		logger.FatalErrLog(err)
	}
}

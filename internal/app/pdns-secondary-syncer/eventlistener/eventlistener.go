package eventlistener

import (
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/config"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/modeljob"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/powerdns"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/primaryzoneprovider"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/worker"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

var (
	changereceivedtotal   = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_received_zone_events_total", ConstLabels: map[string]string{"event_type": "change"}, Help: "The total count of received zone events"})
	deletereceivedtotal   = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_received_zone_events_total", ConstLabels: map[string]string{"event_type": "delete"}, Help: "The total count of received zone events"})
	createreceivedtotal   = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_received_zone_events_total", ConstLabels: map[string]string{"event_type": "create"}, Help: "The total count of received zone events"})
	zonestaterequesttotal = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_zone_state_request_total", Help: "The total count of zone state requests"})
)

func StartEventListenerAndWorker(ms *microservice.Microservice) error {
	serviceconfig := config.GetConfiguration(ms)

	startAddZoneEventListener(ms, serviceconfig)
	startChangeZoneEventListener(ms, serviceconfig)
	startDeleteZoneEventListener(ms, serviceconfig)
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

func startAddZoneEventListener(ms *microservice.Microservice, conf *config.ServiceConfiguration) {
	addtopic := conf.AddEventTopic
	hostaddtopic := addtopic + "." + ms.ID

	addhandler := func(msg *nats.Msg) {
		go worker.EnqueJob(&modeljob.PowerDNSAPIJob{
			Jobtype: modeljob.AddZone,
			Msg:     msg,
			Ms:      ms,
			Conf:    conf,
		})
		createreceivedtotal.Inc()
	}

	ms.MessageBroker.SubscribeAsync(addtopic, addhandler)
	ms.MessageBroker.SubscribeAsync(hostaddtopic, addhandler)
}

func startChangeZoneEventListener(ms *microservice.Microservice, conf *config.ServiceConfiguration) {
	changetopic := conf.ChangeEventTopic
	hostchangetopic := changetopic + "." + ms.ID

	changehandler := func(msg *nats.Msg) {
		go worker.EnqueJob(&modeljob.PowerDNSAPIJob{
			Jobtype: modeljob.ChangeZone,
			Msg:     msg,
			Ms:      ms,
			Conf:    conf,
		})
		changereceivedtotal.Inc()
	}

	ms.MessageBroker.SubscribeAsync(changetopic, changehandler)
	ms.MessageBroker.SubscribeAsync(hostchangetopic, changehandler)
}

func startDeleteZoneEventListener(ms *microservice.Microservice, conf *config.ServiceConfiguration) {
	deltopic := conf.DeleteEventTopic
	hostdeltopic := deltopic + "." + ms.ID

	delhandler := func(msg *nats.Msg) {
		go worker.EnqueJob(&modeljob.PowerDNSAPIJob{
			Jobtype: modeljob.DeleteZone,
			Msg:     msg,
			Ms:      ms,
			Conf:    conf,
		})
		deletereceivedtotal.Inc()
	}

	ms.MessageBroker.SubscribeAsync(deltopic, delhandler)
	ms.MessageBroker.SubscribeAsync(hostdeltopic, delhandler)
}

func startSecondaryZoneStateEventListener(ms *microservice.Microservice, conf *config.ServiceConfiguration) {
	receivertopic := conf.ZoneStateReceiverTopic

	zonestatehandler := func(msg *nats.Msg) {
		go powerdns.RespondSecondaryZoneState(msg, conf)
		logger.DebugLog("[Zone State Request] requested with query: " + string(msg.Data) + " from: " + msg.Reply)
		zonestaterequesttotal.Inc()
	}

	ms.MessageBroker.SubscribeAsync(receivertopic, zonestatehandler)
}

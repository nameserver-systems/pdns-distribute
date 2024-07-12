package eventutils

import (
	"encoding/json"
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelevent"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

func PublishDeleteZoneEvent(ms *microservice.Microservice, topic string, zoneid string) {
	event := modelevent.ZoneDeleteEvent{
		Zone:      zoneid,
		DeletedAt: time.Now(),
	}

	payload, marshalerr := json.Marshal(event)
	if marshalerr != nil {
		logger.ErrorErrLog(marshalerr)
	}

	ms.MessageBroker.PersistedPublish(topic, payload)
	logger.DebugLog("[Delete Zone Event] triggered for zone: " + zoneid + " on topic: " + topic) //nolint:goconst
}

func PublishChangeZoneEvent(ms *microservice.Microservice, topic string, zoneid string) {
	event := modelevent.ZoneChangeEvent{
		Zone:      zoneid,
		ChangedAt: time.Now(),
	}

	payload, marshalerr := json.Marshal(event)
	if marshalerr != nil {
		logger.ErrorErrLog(marshalerr)
	}

	ms.MessageBroker.PersistedPublish(topic, payload)
	logger.DebugLog("[Change Zone Event] triggered for zone: " + zoneid + " on topic: " + topic)
}

func PublishCreateZoneEvent(ms *microservice.Microservice, topic string, zoneid string) {
	event := modelevent.ZoneAddEvent{
		Zone:      zoneid,
		CreatedAt: time.Now(),
	}

	payload, marshalerr := json.Marshal(event)
	if marshalerr != nil {
		logger.ErrorErrLog(marshalerr)
	}

	ms.MessageBroker.PersistedPublish(topic, payload)
	logger.DebugLog("[Create Zone Event] triggered for zone: " + zoneid + " on topic: " + topic)
}

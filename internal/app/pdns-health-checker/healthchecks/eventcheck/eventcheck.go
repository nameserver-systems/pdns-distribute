package eventcheck

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/config"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/dnsutils"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/utils"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/models"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/eventutils"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelevent"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelpowerdns"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/messaging"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	changepublishedtotal = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_event_health_check_published_zone_events_total", ConstLabels: map[string]string{"event_type": "change"}, Help: "The total count of published zone events"})
	deletepublishedtotal = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_event_health_check_published_zone_events_total", ConstLabels: map[string]string{"event_type": "delete"}, Help: "The total count of published zone events"})
	createpublishedtotal = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_event_health_check_published_zone_events_total", ConstLabels: map[string]string{"event_type": "create"}, Help: "The total count of published zone events"})
	changereceivedtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_event_health_check_received_zone_events_total", ConstLabels: map[string]string{"event_type": "change"}, Help: "The total count of received zone events"})
	deletereceivedtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_event_health_check_received_zone_events_total", ConstLabels: map[string]string{"event_type": "delete"}, Help: "The total count of received zone events"})
	createreceivedtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_event_health_check_received_zone_events_total", ConstLabels: map[string]string{"event_type": "create"}, Help: "The total count of received zone events"})
)

func StartEventCheckHandler(ms *microservice.Microservice, conf *config.ServiceConfiguration,
	actualstate *models.State,
) {
	hs := models.InitHealthServiceObject(ms, conf, actualstate)

	waitafterevent := hs.Conf.EventCheckWait

	addevent := hs.Conf.AddEventTopic
	changeevent := hs.Conf.ChangeEventTopic
	deleteevent := hs.Conf.DeleteEventTopic

	consumer := ms.MessageBroker.GetConsumer()
	_, err := consumer.Consume(func(msg jetstream.Msg) {
		subject := msg.Subject()
		if strings.HasPrefix(subject, addevent) {
			go checkAddEvent(msg, hs, waitafterevent)
			createreceivedtotal.Inc()
		} else if strings.HasPrefix(subject, changeevent) {
			go checkChangeEvent(msg, hs, waitafterevent)
			changereceivedtotal.Inc()
		} else if strings.HasPrefix(subject, deleteevent) {
			go checkDeleteEvent(msg, hs, waitafterevent)
			deletereceivedtotal.Inc()
		} else {
			logger.DebugLog("[Zone Event Listener]: not matched on topic: " + subject)
		}
		if err := msg.Ack(); err != nil {
			logger.ErrorErrLog(err)
		}
	})
	if err != nil {
		logger.ErrorErrLog(err)

		return
	}
}

func checkDeleteEvent(msg jetstream.Msg, hs *models.HealthService, waitafterevent time.Duration) {
	data := msg.Data()
	deleteEventObject := modelevent.ZoneDeleteEvent{}

	parseerr := json.Unmarshal(data, &deleteEventObject)
	if parseerr != nil {
		logger.ErrorErrLog(parseerr)
	}

	retrievalerr := hs.CheckFreshnessOfState(deleteEventObject.DeletedAt)
	if retrievalerr != nil {
		logger.ErrorErrLog(retrievalerr)
	}

	time.Sleep(waitafterevent)

	logger.DebugLog("[Delete Event Health Check] start check for zone: " + deleteEventObject.Zone)

	eventcheck := models.EventCheckObject{
		Eventtype: models.DeleteZone,
		Zoneid:    deleteEventObject.Zone,
	}

	checkZoneFreshnessOfSecondaries(eventcheck, hs)
}

func checkChangeEvent(msg jetstream.Msg, hs *models.HealthService, waitafterevent time.Duration) {
	data := msg.Data()
	changeEventObject := modelevent.ZoneChangeEvent{}

	parseerr := json.Unmarshal(data, &changeEventObject)
	if parseerr != nil {
		logger.ErrorErrLog(parseerr)
	}

	retrievalerr := hs.CheckFreshnessOfState(changeEventObject.ChangedAt)
	if retrievalerr != nil {
		logger.ErrorErrLog(retrievalerr)
	}

	time.Sleep(waitafterevent)

	logger.DebugLog("[Change Event Health Check] start check for zone: " + changeEventObject.Zone)

	eventcheck := models.EventCheckObject{
		Eventtype: models.ChangeZone,
		Zoneid:    changeEventObject.Zone,
	}

	checkZoneFreshnessOfSecondaries(eventcheck, hs)
}

func checkAddEvent(msg jetstream.Msg, hs *models.HealthService, waitafterevent time.Duration) {
	data := msg.Data()
	addEventObject := modelevent.ZoneAddEvent{}

	parseerr := json.Unmarshal(data, &addEventObject)
	if parseerr != nil {
		logger.ErrorErrLog(parseerr)
	}

	retrievalerr := hs.CheckFreshnessOfState(addEventObject.CreatedAt)
	if retrievalerr != nil {
		logger.ErrorErrLog(retrievalerr)
	}

	time.Sleep(waitafterevent)

	logger.DebugLog("[Add Event Health Check] start check for zone: " + addEventObject.Zone)

	eventcheck := models.EventCheckObject{
		Eventtype: models.AddZone,
		Zoneid:    addEventObject.Zone,
	}

	checkZoneFreshnessOfSecondaries(eventcheck, hs)
}

func checkZoneFreshnessOfSecondaries(eventcheck models.EventCheckObject, hs *models.HealthService) {
	activesecondaries := hs.State.GetActiveSecondaries()
	primaryzoneserial, retrivalerr := getPrimaryZoneSerial(eventcheck.Zoneid, hs)

	if retrivalerr != nil {
		logger.ErrorErrLog(retrivalerr)
	}

	eventcheck.Primaryserial = primaryzoneserial

	for _, secondary := range activesecondaries {
		checkSecondaryZoneFreshness(eventcheck, secondary, hs)
	}
}

func getPrimaryZoneSerial(zoneid string, hs *models.HealthService) (int32, error) {
	pdnsconnection := modelpowerdns.PDNSconnectionobject{
		PowerDNSurl: hs.Conf.PowerDNSURL,
		ServerID:    hs.Conf.PowerDNSServerID,
		Apitoken:    hs.Conf.PowerDNSAPIToken,
	}

	return dnsutils.GetZoneSerialFromFromPrimary(pdnsconnection, zoneid)
}

func checkSecondaryZoneFreshness(eventcheck models.EventCheckObject, secondary messaging.ResolvedService,
	hs *models.HealthService,
) {
	zoneid := eventcheck.Zoneid
	primaryserial := eventcheck.Primaryserial

	secondaryserial, resolveerr := utils.GetZoneSerialFromSecondary(secondary, hs, zoneid)
	if resolveerr != nil {
		logger.ErrorErrLog(resolveerr)
	}

	logger.DebugLog("[Check Secondary Freshness] check for zone: " + zoneid + " primary serial: " +
		strconv.Itoa(int(primaryserial)) + " secondary serial: " + strconv.Itoa(int(secondaryserial)) + " secondary: " + secondary.ID)

	if primaryserial != secondaryserial {
		switch eventcheck.Eventtype {
		case models.AddZone:
			eventutils.PublishCreateZoneEvent(hs.Ms, utils.AppendIDToTopic(hs.Conf.AddEventTopic, secondary), zoneid)
			createpublishedtotal.Inc()
		case models.DeleteZone:
			eventutils.PublishDeleteZoneEvent(hs.Ms, utils.AppendIDToTopic(hs.Conf.DeleteEventTopic, secondary), zoneid)
			deletepublishedtotal.Inc()
		case models.ChangeZone:
			eventutils.PublishChangeZoneEvent(hs.Ms, utils.AppendIDToTopic(hs.Conf.ChangeEventTopic, secondary), zoneid)
			changepublishedtotal.Inc()
		}
	}
}

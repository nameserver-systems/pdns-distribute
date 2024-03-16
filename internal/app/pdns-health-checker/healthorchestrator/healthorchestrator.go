package healthorchestrator

import (
	"fmt"
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/config"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/dnsutils"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/models"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/servicediscovery"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelpowerdns"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	primaryzonecount     = promauto.NewGauge(prometheus.GaugeOpts{Name: "healthchecker_actual_primary_zone_count", Help: "The actual count of zones on primary"})
	activesecondarycount = promauto.NewGauge(prometheus.GaugeOpts{Name: "healthchecker_actual_active_secondary_count", Help: "The actual count of active secondaries"})
	refreshcyclestotal   = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_refresh_cycles_total", Help: "The total count of refresh cycles of primary zones and active secondaries"})
)

func StartHealthServices(ms *microservice.Microservice) error {
	serviceconfig := config.GetConfiguration(ms)

	port, err := ms.GetServicePort()
	if err != nil {
		return err
	}

	stream, err := ms.MessageBroker.CreatePersistentMessageStore("pdns-distribute-event-store", []string{"zone.>"})
	if err != nil {
		return err
	}

	consumer, err := ms.MessageBroker.CreatePersistentMessageReceiver("health-checker-event-client", ms.ID, "add", port, "healthchecker", stream)
	if err != nil {
		return err
	}

	ms.MessageBroker.SetStream(stream)
	ms.MessageBroker.SetConsumer(consumer)

	actualstate := models.GenerateStateObject()

	err = SetActiveZonesInState(serviceconfig, actualstate)
	if err != nil {
		return err
	}

	err = SetActiveSecondariesInState(ms, actualstate)
	if err != nil {
		return err
	}

	go startActiveZoneAndSecondaryRefresh(ms, serviceconfig, actualstate)
	go healthchecks.StartAllHealthChecks(ms, serviceconfig, actualstate)

	return nil
}

func SetActiveZonesInState(serviceconfig *config.ServiceConfiguration, actualstate *models.State) error {
	pdnsconnection := modelpowerdns.PDNSconnectionobject{
		PowerDNSurl: serviceconfig.PowerDNSURL,
		ServerID:    serviceconfig.PowerDNSServerID,
		Apitoken:    serviceconfig.PowerDNSAPIToken,
	}

	expectedzones, retrievalerr := dnsutils.GetActiveZonesFromPrimary(pdnsconnection)
	if retrievalerr != nil {
		return retrievalerr
	}

	actualstate.SetExpectedZoneMap(expectedzones)

	logger.DebugLog("[Set Active Primary Zones] zone state: " + fmt.Sprint(expectedzones))

	primaryzonecount.Set(float64(len(expectedzones)))

	return nil
}

func SetActiveSecondariesInState(ms *microservice.Microservice, actualstate *models.State) error {
	activesecondaries, resolveerr := servicediscovery.GetActiveSecondaries(ms)
	if resolveerr != nil {
		return resolveerr
	}

	actualstate.SetActiveSecondaries(activesecondaries)

	logger.DebugLog("[Set Active Secondaries] secondary state: " + fmt.Sprint(activesecondaries))

	activesecondarycount.Set(float64(len(activesecondaries)))

	return nil
}

func startActiveZoneAndSecondaryRefresh(ms *microservice.Microservice, conf *config.ServiceConfiguration, actualstate *models.State) { //nolint:lll
	refreshintervall := conf.ActiveZoneSecondaryRefreshIntervall
	refreshsignal := time.Tick(refreshintervall) //nolint:staticcheck

	go func() {
		for range refreshsignal {
			logger.DebugLog("[Zone/Secondary Refresh]: start refresh cycle")
			refreshSecondaries(ms, actualstate)
			refreshZones(conf, actualstate)
			refreshcyclestotal.Inc()
		}
	}()
}

func refreshSecondaries(ms *microservice.Microservice, actualstate *models.State) {
	err := SetActiveSecondariesInState(ms, actualstate)
	if err != nil {
		logger.ErrorErrLog(err)
	}
}

func refreshZones(conf *config.ServiceConfiguration, actualstate *models.State) {
	err := SetActiveZonesInState(conf, actualstate)
	if err != nil {
		logger.ErrorErrLog(err)
	}
}

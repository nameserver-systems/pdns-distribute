package intervallsigningsync

import (
	"fmt"
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/config"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/dnsutils"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/utils"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/models"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/eventutils"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelpowerdns"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelzone"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/servicediscovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	signingsynctotal       = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_singing_sync_cycles_total", Help: "The total count of signing sync cycles of dnssec zones"})
	Dnssecprimaryzonecount = promauto.NewGauge(prometheus.GaugeOpts{Name: "healthchecker_actual_dnssec_primary_zone_count", Help: "The actual count of dnssec zones on primary"})
	changepublishedtotal   = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_interval_signing_sync_published_zone_events_total", ConstLabels: map[string]string{"event_type": "change"}, Help: "The total count of published zone events"})
)

func StartPeridoicalSigningSync(ms *microservice.Microservice, conf *config.ServiceConfiguration,
	actualstate *models.State,
) {
	const refreshintervall = 7 * 24 * time.Hour

	hs := models.InitHealthServiceObject(ms, conf, actualstate)

	waitToStart := getInitialWaitTimeForFirstSignatureRefresh()

	logger.DebugLog("[Periodical Signing Sync] start waiting for initial sync")

	time.Sleep(waitToStart)

	go refreshSignaturesOnSecondaries(hs)

	triggersync := time.Tick(refreshintervall) //nolint:staticcheck

	go func() {
		for range triggersync {
			refreshSignaturesOnSecondaries(hs)
		}
	}()
}

func getInitialWaitTimeForFirstSignatureRefresh() time.Duration {
	const daysinweek = 7

	now := time.Now().UTC()

	diff := time.Thursday - now.Weekday()
	if diff < 0 {
		diff = daysinweek + diff
	}

	future := now.AddDate(0, 0, int(diff)).Add(time.Hour)
	waitToStart := future.Sub(now)

	return waitToStart
}

func refreshSignaturesOnSecondaries(hs *models.HealthService) {
	activesecondaries := hs.State.GetActiveSecondaries()

	pdnsconnection := modelpowerdns.PDNSconnectionobject{
		PowerDNSurl: hs.Conf.PowerDNSURL,
		ServerID:    hs.Conf.PowerDNSServerID,
		Apitoken:    hs.Conf.PowerDNSAPIToken,
	}

	activeprimaryzones, err := dnsutils.GetActiveDNSSECZonesFromPrimary(pdnsconnection)
	if err != nil {
		logger.ErrorErrLog(err)
	}

	logger.DebugLog("[Periodical Signing Sync] start sync for secondaries: " + fmt.Sprint(activesecondaries) + " primary zones: " +
		fmt.Sprint(activeprimaryzones))

	for _, secondary := range activesecondaries {
		refreshSecondarySignatures(secondary, activeprimaryzones, hs)
	}

	Dnssecprimaryzonecount.Set(float64(len(activeprimaryzones)))
	signingsynctotal.Inc()
}

func refreshSecondarySignatures(secondary servicediscovery.ResolvedService, primaryzonestatemap modelzone.Zonestatemap,
	hs *models.HealthService) {
	for zoneid := range primaryzonestatemap {
		eventutils.PublishChangeZoneEvent(hs.Ms, utils.AppendIDToTopic(hs.Conf.ChangeEventTopic, secondary), zoneid)
		changepublishedtotal.Inc()
	}
}

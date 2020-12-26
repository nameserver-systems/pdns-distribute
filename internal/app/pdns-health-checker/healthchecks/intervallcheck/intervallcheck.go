package intervallcheck

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/config"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/utils"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/models"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/eventutils"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/modelzone"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/servicediscovery"
)

var (
	intervalcyclestotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_interval_health_cycles_total", Help: "The total count of interval health check cycles"})
	changepublishedtotal = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_interval_health_published_zone_events_total", ConstLabels: map[string]string{"event_type": "change"}, Help: "The total count of published zone events"})
	deletepublishedtotal = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_interval_health_published_zone_events_total", ConstLabels: map[string]string{"event_type": "delete"}, Help: "The total count of published zone events"})
	createpublishedtotal = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_interval_health_published_zone_events_total", ConstLabels: map[string]string{"event_type": "create"}, Help: "The total count of published zone events"})
)

func StartPeriodicalCheck(ms *microservice.Microservice, conf *config.ServiceConfiguration, actualstate *models.State) {
	hs := models.InitHealthServiceObject(ms, conf, actualstate)

	checkintervall := conf.PeriodicalCheckIntervall
	triggersync := time.Tick(checkintervall) //nolint:staticcheck

	go func() {
		for range triggersync {
			checkFreshnessOfSecondaries(hs)
		}
	}()
}

func checkFreshnessOfSecondaries(hs *models.HealthService) {
	activesecondaries := hs.State.GetActiveSecondaries()
	activeprimaryzones := hs.State.GetExpectedZoneMap()

	logger.DebugLog("[Interval Health Check] start check for secondaries: " + fmt.Sprint(activesecondaries) + " primary zones: " +
		fmt.Sprint(activeprimaryzones))

	for _, secondaries := range activesecondaries {
		err := checkSecondaryFreshness(secondaries, activeprimaryzones, hs)
		if err != nil {
			logger.ErrorErrLog(err)
		}
	}

	intervalcyclestotal.Inc()
}

func checkSecondaryFreshness(secondary servicediscovery.ResolvedService, primaryzonestatemap modelzone.Zonestatemap,
	hs *models.HealthService) error {
	secondaryzonestatemap, err := utils.GetSecondaryZoneStateMap(secondary, hs)
	if err != nil {
		return err
	}

	addedzones, deletedzones, changedzones := primaryzonestatemap.Diff(secondaryzonestatemap)

	for zoneid := range addedzones {
		eventutils.PublishCreateZoneEvent(hs.Ms, utils.AppendIDToTopic(hs.Conf.AddEventTopic, secondary), zoneid)
		createpublishedtotal.Inc()
	}

	for zoneid := range deletedzones {
		eventutils.PublishDeleteZoneEvent(hs.Ms, utils.AppendIDToTopic(hs.Conf.DeleteEventTopic, secondary), zoneid)
		deletepublishedtotal.Inc()
	}

	for zoneid := range changedzones {
		eventutils.PublishChangeZoneEvent(hs.Ms, utils.AppendIDToTopic(hs.Conf.ChangeEventTopic, secondary), zoneid)
		changepublishedtotal.Inc()
	}

	return nil
}

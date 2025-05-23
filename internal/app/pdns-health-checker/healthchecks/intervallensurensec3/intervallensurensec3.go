package intervallensurensec3

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/config"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/dnsutils"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/intervallsigningsync"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/models"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/eventutils"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelpowerdns"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const expectedNSEC3PARAMS = "1 0 0 -" // RFC9276

var (
	ensurensec3cycles    = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_ensure_nsec3_cyles_total", Help: "The total count of ensure nse3 cycles for dnssec zones"})
	powerdnsclierrtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_powerdns_cli_call_total", ConstLabels: map[string]string{"state": "failed"}, Help: "The total count of powerdns cli calls"})
	powerdnsclicalltotal = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_powerdns_cli_call_total", ConstLabels: map[string]string{"state": "successful"}, Help: "The total count of powerdns cli calls"})
)

func StartIntervallEnsureNsec3(ms *microservice.Microservice, conf *config.ServiceConfiguration,
	actualstate *models.State,
) {
	refreshintervall := conf.NSEC3CheckIntervall

	hs := models.InitHealthServiceObject(ms, conf, actualstate)

	logger.DebugLog("[Interval Ensure Nsec3] start waiting for initial sync")

	go ensurensec3(hs)

	triggersync := time.Tick(refreshintervall) //nolint:staticcheck

	go func() {
		for range triggersync {
			ensurensec3(hs)
			ensurensec3cycles.Inc()
		}
	}()
}

func ensurensec3(hs *models.HealthService) {
	pdnsconnection := modelpowerdns.PDNSconnectionobject{
		PowerDNSurl: hs.Conf.PowerDNSURL,
		ServerID:    hs.Conf.PowerDNSServerID,
		Apitoken:    hs.Conf.PowerDNSAPIToken,
	}

	activeprimaryzones, err := dnsutils.GetActiveDNSSECZonesFromPrimary(pdnsconnection)
	if err != nil {
		logger.ErrorErrLog(err)
	}

	logger.DebugLog("[Interval Ensure Nsec3] start setting nsec3 for zones: " + fmt.Sprint(activeprimaryzones))

	for zone := range activeprimaryzones {
		checkNecessityForUpdateNSEC3(hs, pdnsconnection, zone)
	}

	intervallsigningsync.Dnssecprimaryzonecount.Set(float64(len(activeprimaryzones)))
}

func checkNecessityForUpdateNSEC3(hs *models.HealthService, pdnsconnection modelpowerdns.PDNSconnectionobject, zoneID string) {
	mData, mdErr := dnsutils.GetZoneMetaDataFromPrimary(pdnsconnection, zoneID, "NSEC3PARAM")
	if mdErr != nil {
		logger.ErrorErrLog(mdErr)
	}

	if len(mData) > 0 {
		if mData[0] == expectedNSEC3PARAMS {
			return
		}
	}
	err := setnsec3(zoneID)
	if err != nil {
		logger.ErrorErrLog(err)
	}

	rectifyErr := dnsutils.RectifyZone(pdnsconnection, zoneID)
	if rectifyErr != nil {
		logger.ErrorErrLog(rectifyErr)
	}

	eventutils.PublishChangeZoneEvent(hs.Ms, hs.Conf.ChangeEventTopic, zoneID)
}

func setnsec3(zoneid string) error {
	metaout, pdnsutilmetaerr := exec.Command("pdnsutil", "set-nsec3", zoneid).Output()
	if pdnsutilmetaerr != nil {
		logger.ErrorLog(string(metaout))
		powerdnsclierrtotal.Inc()

		return pdnsutilmetaerr
	}

	powerdnsclicalltotal.Inc()

	return nil
}

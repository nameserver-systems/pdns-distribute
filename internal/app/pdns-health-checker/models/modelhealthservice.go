package models

import (
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/config"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/dnsutils"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelpowerdns"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
)

type HealthService struct {
	PDNSConnection modelpowerdns.PDNSconnectionobject
	State          *State

	Ms   *microservice.Microservice
	Conf *config.ServiceConfiguration
}

func InitHealthServiceObject(ms *microservice.Microservice, conf *config.ServiceConfiguration,
	state *State) *HealthService {
	hs := &HealthService{
		State: state,
		Ms:    ms,
		Conf:  conf,
	}
	hs.initPDNSConnectionObject()

	return hs
}

func (hs *HealthService) initPDNSConnectionObject() {
	hs.PDNSConnection = modelpowerdns.PDNSconnectionobject{
		PowerDNSurl: hs.Conf.PowerDNSURL,
		ServerID:    hs.Conf.PowerDNSServerID,
		Apitoken:    hs.Conf.PowerDNSAPIToken,
	}
}

func (hs *HealthService) CheckFreshnessOfState(msgtime time.Time) error {
	if msgtime.After(hs.State.GetExpectedZoneChangeTime()) {
		expectedzones, retrievalerr := dnsutils.GetActiveZonesFromPrimary(hs.PDNSConnection)
		if retrievalerr != nil {
			return retrievalerr
		}

		hs.State.SetExpectedZoneMap(expectedzones)
	}

	return nil
}

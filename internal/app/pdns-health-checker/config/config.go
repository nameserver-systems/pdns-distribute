package config

import (
	"time"

	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice"
)

type ServiceConfiguration struct {
	AddEventTopic                       string
	ChangeEventTopic                    string
	DeleteEventTopic                    string
	PowerDNSURL                         string
	PowerDNSAPIToken                    string
	PowerDNSServerID                    string
	EventCheckWait                      time.Duration
	ActiveZoneSecondaryRefreshIntervall time.Duration
	PeriodicalCheckIntervall            time.Duration
	NSEC3CheckIntervall                 time.Duration
	ZoneStateReceiverTopic              string
}

func GetConfiguration(ms *microservice.Microservice) *ServiceConfiguration {
	return &ServiceConfiguration{
		AddEventTopic:                       ms.Config.GetStringSetting("ZoneEventTopics.Add"),
		ChangeEventTopic:                    ms.Config.GetStringSetting("ZoneEventTopics.Mod"),
		DeleteEventTopic:                    ms.Config.GetStringSetting("ZoneEventTopics.Del"),
		PowerDNSURL:                         ms.Config.GetStringSetting("PowerDNS.URL"),
		PowerDNSAPIToken:                    ms.Config.GetStringSetting("PowerDNS.APIToken"),
		PowerDNSServerID:                    ms.Config.GetStringSetting("PowerDNS.ServerID"),
		EventCheckWait:                      ms.Config.GetTimeDuration("HealthChecks.EventCheckWaitTime"),
		ActiveZoneSecondaryRefreshIntervall: ms.Config.GetTimeDuration("HealthChecks.ActiveZoneSecondaryRefreshIntervall"),
		PeriodicalCheckIntervall:            ms.Config.GetTimeDuration("HealthChecks.PeriodicalCheckIntervall"),
		NSEC3CheckIntervall:                 ms.Config.GetTimeDuration("HealthChecks.NSEC3CheckIntervall"),
		ZoneStateReceiverTopic:              ms.Config.GetStringSetting("ZoneStateTopics.Prefix"),
	}
}

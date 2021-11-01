package config

import (
	"time"

	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
)

type ServiceConfiguration struct {
	ReceiverTopic       string
	PowerDNSURL         string
	PowerDNSAPIToken    string
	PowerDNSServerID    string
	PowerDNSAXFRAddress string
	PowerDNSAXFRTimeout time.Duration
}

func GetConfiguration(ms *microservice.Microservice) *ServiceConfiguration {
	return &ServiceConfiguration{
		ReceiverTopic:       ms.Config.GetStringSetting("ZoneDataTopics.Wildcard"),
		PowerDNSURL:         ms.Config.GetStringSetting("PowerDNS.URL"),
		PowerDNSAPIToken:    ms.Config.GetStringSetting("PowerDNS.APIToken"),
		PowerDNSServerID:    ms.Config.GetStringSetting("PowerDNS.ServerID"),
		PowerDNSAXFRTimeout: ms.Config.GetTimeDuration("PowerDNS.AXFRTimeout"),
		PowerDNSAXFRAddress: ms.Config.GetStringSetting("PowerDNS.AXFRAddress"),
	}
}

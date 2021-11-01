package config

import "github.com/nameserver-systems/pdns-distribute/pkg/microservice"

type ServiceConfiguration struct {
	AddEventTopic    string
	ChangeEventTopic string
	DeleteEventTopic string
	PowerDNSURL      string
	ServiceURL       string
	CertPath         string
	KeyPath          string
}

func GetConfiguration(ms *microservice.Microservice) *ServiceConfiguration {
	return &ServiceConfiguration{
		AddEventTopic:    ms.Config.GetStringSetting("ZoneEventTopics.Add"),
		ChangeEventTopic: ms.Config.GetStringSetting("ZoneEventTopics.Mod"),
		DeleteEventTopic: ms.Config.GetStringSetting("ZoneEventTopics.Del"),
		PowerDNSURL:      ms.Config.GetStringSetting("PowerDNS.URL"),
		ServiceURL:       ms.Config.GetStringSetting("Service.URL"),
		CertPath:         ms.Config.GetStringSetting("Service.Cert"),
		KeyPath:          ms.Config.GetStringSetting("Service.Key"),
	}
}

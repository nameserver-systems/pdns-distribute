package config

import (
	"time"

	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
)

type ServiceConfiguration struct {
	ServiceURL             string
	AddEventTopic          string
	ChangeEventTopic       string
	DeleteEventTopic       string
	SenderTopic            string
	PowerDNSURL            string
	PowerDNSAPIToken       string
	PowerDNSServerID       string
	EventDelay             time.Duration
	APIWorker              int
	ZoneStateReceiverTopic string
	AXFRPrimaryAddress     string
}

func GetConfiguration(ms *microservice.Microservice) *ServiceConfiguration {
	sendertopic := generateSenderTopic(ms)
	zonestatereceivertopic := generateZoneStateReceiverTopic(ms)

	return &ServiceConfiguration{
		ServiceURL:             ms.Config.GetStringSetting("Service.URL"),
		AddEventTopic:          ms.Config.GetStringSetting("ZoneEventTopics.Add"),
		ChangeEventTopic:       ms.Config.GetStringSetting("ZoneEventTopics.Mod"),
		DeleteEventTopic:       ms.Config.GetStringSetting("ZoneEventTopics.Del"),
		SenderTopic:            sendertopic,
		PowerDNSURL:            ms.Config.GetStringSetting("PowerDNS.URL"),
		PowerDNSAPIToken:       ms.Config.GetStringSetting("PowerDNS.APIToken"),
		PowerDNSServerID:       ms.Config.GetStringSetting("PowerDNS.ServerID"),
		EventDelay:             ms.Config.GetTimeDuration("PowerDNS.EventDelay"),
		APIWorker:              ms.Config.GetIntSetting("PowerDNS.APIWorker"),
		ZoneStateReceiverTopic: zonestatereceivertopic,
		AXFRPrimaryAddress:     ms.Config.GetStringSetting("AXFRPrimary.Address"),
	}
}

func generateSenderTopic(ms *microservice.Microservice) string {
	sendertopic := ms.Config.GetStringSetting("ZoneDataTopics.Prefix") + ms.ID

	return sendertopic
}

func generateZoneStateReceiverTopic(ms *microservice.Microservice) string {
	sendertopic := ms.Config.GetStringSetting("ZoneStateTopics.Prefix") + ms.ID

	return sendertopic
}

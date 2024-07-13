package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/models"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelzone"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/servicediscovery"
)

func AppendIDToTopic(topic string, secondary servicediscovery.ResolvedService) string {
	return topic + "." + secondary.ID
}

func GetSecondaryZoneStateMap(secondary servicediscovery.ResolvedService,
	hs *models.HealthService) (modelzone.Zonestatemap, error) {
	const messagetimeout = 10 * time.Second

	var transferobject modelzone.Zonestatemaptransferobject

	zonestatetopic := hs.Conf.ZoneStateReceiverTopic + secondary.ID

	zonestatemaprequest := modelzone.Zonestaterequestevent{
		RequestedAt: time.Now(),
	}

	paylaod, marshalerr := json.Marshal(zonestatemaprequest)
	if marshalerr != nil {
		return nil, marshalerr
	}

	response, retrieveerr := hs.Ms.MessageBroker.PublishRequestAndWait(zonestatetopic, paylaod, messagetimeout)
	if retrieveerr != nil {
		return nil, retrieveerr
	}

	responsedata := response.Data

	unmarshalerr := json.Unmarshal(responsedata, &transferobject)
	if unmarshalerr != nil {
		return nil, unmarshalerr
	}

	logger.DebugLog("[Get Secondary Zone State] secondary topic: " + zonestatetopic + " zone state: " +
		fmt.Sprint(transferobject.Statemap))

	return transferobject.Statemap, nil
}

func GetZoneSerialFromSecondary(secondary servicediscovery.ResolvedService, hs *models.HealthService,
	zoneid string) (int32, error) {
	const expectedzonecount = 1

	const messagetimeout = 10 * time.Second

	var transferobject modelzone.Zonestatemaptransferobject

	zonestatetopic := hs.Conf.ZoneStateReceiverTopic + secondary.ID

	zonestaterequest := modelzone.Zonestaterequestevent{
		Zone:        zoneid,
		RequestedAt: time.Now(),
	}

	paylaod, marshalerr := json.Marshal(zonestaterequest)
	if marshalerr != nil {
		return int32(0), marshalerr
	}

	response, retrieveerr := hs.Ms.MessageBroker.PublishRequestAndWait(zonestatetopic, paylaod, messagetimeout)
	if retrieveerr != nil {
		return int32(0), retrieveerr
	}

	responsedata := response.Data

	unmarshalerr := json.Unmarshal(responsedata, &transferobject)
	if unmarshalerr != nil {
		return int32(0), unmarshalerr
	}

	zone := transferobject.Statemap

	if len(zone) != expectedzonecount {
		return int32(0), errUnexpectedZoneCount
	}

	if _, exists := zone[zoneid]; !exists {
		return int32(0), errZoneIDNotInStateMap
	}

	return zone[zoneid], nil
}

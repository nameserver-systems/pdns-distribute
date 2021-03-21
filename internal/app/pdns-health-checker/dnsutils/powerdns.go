package dnsutils

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/modelpowerdns"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/modelzone"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/powerdnsutils"
)

var (
	powerdnsapicalltotal = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_powerdns_api_call_total", ConstLabels: map[string]string{"state": "successful"}, Help: "The total count of powerdns api calls"})
	powerdnsapierrtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "healthchecker_powerdns_api_call_total", ConstLabels: map[string]string{"state": "failed"}, Help: "The total count of powerdns api calls"})
)

func GetActiveZonesFromPrimary(pdnsconnection modelpowerdns.PDNSconnectionobject) (modelzone.Zonestatemap, error) {
	var actualzones []modelpowerdns.Zone

	var expectedzones modelzone.Zonestatemap

	response, err := powerdnsutils.GetZoneList(pdnsconnection, "", false)
	if err != nil {
		powerdnsapierrtotal.Inc()

		return nil, err
	}

	powerdnsapicalltotal.Inc()

	parseerr := json.Unmarshal([]byte(response), &actualzones)
	if parseerr != nil {
		return nil, parseerr
	}

	expectedzones = powerdnsutils.TransferPowerDNSZonesIntoZoneStateMap(actualzones)

	return expectedzones, nil
}

func GetActiveDNSSECZonesFromPrimary(pdnsconnection modelpowerdns.PDNSconnectionobject) (modelzone.Zonestatemap, error) {
	var actualzones []modelpowerdns.Zone

	var expectedzones modelzone.Zonestatemap

	response, err := powerdnsutils.GetZoneList(pdnsconnection, "", true)
	if err != nil {
		powerdnsapierrtotal.Inc()

		return nil, err
	}

	powerdnsapicalltotal.Inc()

	parseerr := json.Unmarshal([]byte(response), &actualzones)
	if parseerr != nil {
		return nil, parseerr
	}

	expectedzones = powerdnsutils.TransferPowerDNSDNSSECZonesIntoZoneStateMap(actualzones)

	return expectedzones, nil
}

func GetZoneMetaDataFromPrimary(pdnsconnection modelpowerdns.PDNSconnectionobject, zoneid, metadatakind string) ([]string, error) {
	var metadata modelpowerdns.Metadata

	response, err := powerdnsutils.GetZoneMetaData(pdnsconnection, zoneid, metadatakind)
	if err != nil {
		powerdnsapierrtotal.Inc()

		return nil, err
	}

	powerdnsapicalltotal.Inc()

	parseerr := json.Unmarshal([]byte(response), &metadata)
	if parseerr != nil {
		return []string{""}, parseerr
	}

	return metadata.Metadata, nil
}

func RectifyZone(pdnsconnection modelpowerdns.PDNSconnectionobject, zoneid string) error {
	err := powerdnsutils.RectifyZone(pdnsconnection, zoneid)
	if err != nil {
		powerdnsapierrtotal.Inc()

		return err
	}

	powerdnsapicalltotal.Inc()

	return nil
}

func GetZoneSerialFromFromPrimary(pdnsconnection modelpowerdns.PDNSconnectionobject, zoneid string) (int32, error) {
	const expectedzonecount = 1

	var zone []modelpowerdns.Zone

	response, err := powerdnsutils.GetZoneList(pdnsconnection, zoneid, false)
	if err != nil {
		powerdnsapierrtotal.Inc()

		return int32(0), err
	}

	powerdnsapicalltotal.Inc()

	if parseerr := json.Unmarshal([]byte(response), &zone); parseerr != nil {
		return int32(0), parseerr
	}

	if len(zone) != expectedzonecount {
		return int32(0), errTooMuchZones
	}

	if zone[0].ID != zoneid {
		return int32(0), errZoneIDMismatch
	}

	return zone[0].Serial, nil
}

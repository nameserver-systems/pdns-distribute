package powerdnsutils

import (
	"bytes"
	"net/http"
	"strconv"

	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/httputils"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/modelpowerdns"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/modelzone"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

func TransferPowerDNSZonesIntoZoneStateMap(actualzones []modelpowerdns.Zone) modelzone.Zonestatemap {
	expectedzones := make(modelzone.Zonestatemap)

	for _, zone := range actualzones {
		expectedzones[zone.ID] = zone.Serial
	}

	return expectedzones
}

func TransferPowerDNSDNSSECZonesIntoZoneStateMap(actualzones []modelpowerdns.Zone) modelzone.Zonestatemap {
	expectedzones := make(modelzone.Zonestatemap)

	for _, zone := range actualzones {
		if zone.Dnssec {
			expectedzones[zone.ID] = zone.Serial
		}
	}

	return expectedzones
}

func GetZoneList(con modelpowerdns.PDNSconnectionobject, zoneid string, dnssecinfo bool) (string, error) {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + "/api/v1/servers/" + serverid + "/zones?rrsets=false&dnssec=" +
		strconv.FormatBool(dnssecinfo)

	if zoneid != "" {
		zoneDataFileURL += "&zone=" + zoneid
	}

	response, requesterr := httputils.ExecutePowerDNSRequest(http.MethodGet, zoneDataFileURL, apitoken, nil)
	if requesterr != nil {
		return "", requesterr
	}

	return response, nil
}

func DeleteZone(con modelpowerdns.PDNSconnectionobject, zoneid string) error {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + "/api/v1/servers/" + serverid + "/zones/" + zoneid

	_, delerr := httputils.ExecutePowerDNSRequest(http.MethodDelete, zoneDataFileURL, apitoken, nil)
	if delerr != nil {
		return delerr
	}

	return nil
}

func CreateZone(con modelpowerdns.PDNSconnectionobject, storepayload []byte) error {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + "/api/v1/servers/" + serverid + "/zones?rrsets=false"

	// T8D8
	_, storeerr := httputils.ExecutePowerDNSRequest(http.MethodPost, zoneDataFileURL, apitoken,
		bytes.NewBuffer(storepayload))
	if storeerr != nil {
		return storeerr
	}

	return nil
}

func ClearCache(zoneid string, con modelpowerdns.PDNSconnectionobject) {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + "/api/v1/servers/" + serverid + "/cache/flush?domain=" + zoneid

	_, clearerr := httputils.ExecutePowerDNSRequest(http.MethodPut, zoneDataFileURL, apitoken, nil)
	if clearerr != nil {
		logger.ErrorErrLog(clearerr)
	}
}

func GetZoneMetaData(con modelpowerdns.PDNSconnectionobject, zoneid, metadatakind string) (string, error) {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + "/api/v1/servers/" + serverid + "/zones/" + zoneid + "/metadata/" + metadatakind

	resp, storeerr := httputils.ExecutePowerDNSRequest(http.MethodGet, zoneDataFileURL, apitoken, nil)
	if storeerr != nil {
		return "", storeerr
	}

	return resp, nil
}

func RectifyZone(con modelpowerdns.PDNSconnectionobject, zoneid string) error {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + "/api/v1/servers/" + serverid + "/zones/" + zoneid + "/rectify"

	_, storeerr := httputils.ExecutePowerDNSRequest(http.MethodPut, zoneDataFileURL, apitoken, nil)
	if storeerr != nil {
		return storeerr
	}

	return nil
}

func AXFRRetrieve(con modelpowerdns.PDNSconnectionobject, zoneid string) error {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + "/api/v1/servers/" + serverid + "/zones/" + zoneid + "/axfr-retrieve"

	_, storeerr := httputils.ExecutePowerDNSRequest(http.MethodPut, zoneDataFileURL, apitoken, nil)
	if storeerr != nil {
		return storeerr
	}

	return nil
}

func SetZoneMetaData(con modelpowerdns.PDNSconnectionobject, zoneid string, storepayload []byte) error {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + "/api/v1/servers/" + serverid + "/zones/" + zoneid + "/metadata"

	_, storeerr := httputils.ExecutePowerDNSRequest(http.MethodPost, zoneDataFileURL, apitoken,
		bytes.NewBuffer(storepayload))
	if storeerr != nil {
		return storeerr
	}

	return nil
}

package powerdnsutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nameserver-systems/pdns-distribute/internal/pkg/httputils"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelpowerdns"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelzone"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
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
	zoneDataFileURL := con.PowerDNSurl + PowerDNSServerBaseURL + serverid + "/zones?rrsets=false&dnssec=" +
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

func DoesZoneExist(con modelpowerdns.PDNSconnectionobject, zoneID string) (bool, error) {
	serverID := con.ServerID
	apiToken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + PowerDNSServerBaseURL + serverID + "/zones?zone=" + zoneID

	response, requesterr := httputils.ExecutePowerDNSRequest(http.MethodGet, zoneDataFileURL, apiToken, nil)
	if requesterr != nil {
		return false, requesterr
	}

	var zoneList []modelpowerdns.Zone

	parseerr := json.Unmarshal([]byte(response), &zoneList)
	if parseerr != nil {
		return false, parseerr
	}

	if len(zoneList) != 1 {
		return false, nil
	}

	return true, nil
}

func DeleteZone(con modelpowerdns.PDNSconnectionobject, zoneid string) error {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + PowerDNSServerBaseURL + serverid + PowerDNSZoneURLPath + zoneid

	_, delerr := httputils.ExecutePowerDNSRequest(http.MethodDelete, zoneDataFileURL, apitoken, nil)
	if delerr != nil {
		return delerr
	}

	return nil
}

func CreateZone(con modelpowerdns.PDNSconnectionobject, storepayload []byte) error {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + PowerDNSServerBaseURL + serverid + "/zones?rrsets=false"

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
	zoneDataFileURL := con.PowerDNSurl + PowerDNSServerBaseURL + serverid + "/cache/flush?domain=" + zoneid

	_, clearerr := httputils.ExecutePowerDNSRequest(http.MethodPut, zoneDataFileURL, apitoken, nil)
	if clearerr != nil {
		logger.ErrorErrLog(clearerr)
	}
}

func GetZoneMetaData(con modelpowerdns.PDNSconnectionobject, zoneid, metadatakind string) (string, error) {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + PowerDNSServerBaseURL + serverid + PowerDNSZoneURLPath + zoneid + "/metadata/" + metadatakind

	resp, storeerr := httputils.ExecutePowerDNSRequest(http.MethodGet, zoneDataFileURL, apitoken, nil)
	if storeerr != nil {
		return "", storeerr
	}

	return resp, nil
}

func RectifyZone(con modelpowerdns.PDNSconnectionobject, zoneid string) error {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + PowerDNSServerBaseURL + serverid + PowerDNSZoneURLPath + zoneid + "/rectify"

	_, storeerr := httputils.ExecutePowerDNSRequest(http.MethodPut, zoneDataFileURL, apitoken, nil)
	if storeerr != nil {
		return storeerr
	}

	return nil
}

func AXFRRetrieve(con modelpowerdns.PDNSconnectionobject, zoneid string) error {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + PowerDNSServerBaseURL + serverid + PowerDNSZoneURLPath + zoneid + "/axfr-retrieve"

	_, storeerr := httputils.ExecutePowerDNSRequest(http.MethodPut, zoneDataFileURL, apitoken, nil)
	if storeerr != nil {
		return storeerr
	}

	return nil
}

func SetZoneMetaData(con modelpowerdns.PDNSconnectionobject, zoneid string, storepayload []byte) error {
	serverid := con.ServerID
	apitoken := con.Apitoken
	zoneDataFileURL := con.PowerDNSurl + PowerDNSServerBaseURL + serverid + PowerDNSZoneURLPath + zoneid + "/metadata"

	_, storeerr := httputils.ExecutePowerDNSRequest(http.MethodPost, zoneDataFileURL, apitoken,
		bytes.NewBuffer(storepayload))
	if storeerr != nil {
		return storeerr
	}

	return nil
}

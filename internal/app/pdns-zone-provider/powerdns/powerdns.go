//nolint:lll
package powerdns

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-zone-provider/config"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/httputils"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/modelevent"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/modelpowerdns"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/powerdnsutils"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

var (
	powerdnsapicalltotal    = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_powerdns_api_call_total", ConstLabels: map[string]string{"state": "successful"}, Help: "The total count of powerdns api calls"})
	powerdnsapierrtotal     = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_powerdns_api_call_total", ConstLabels: map[string]string{"state": "failed"}, Help: "The total count of powerdns api calls"})
	powerdnsdnstcpcalltotal = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_powerdns_dns_call_total", ConstLabels: map[string]string{"state": "successful", "protocol": "tcp"}, Help: "The total count of powerdns dns calls"})
	powerdnsdnstcperrtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_powerdns_dns_call_total", ConstLabels: map[string]string{"state": "failed", "protocol": "tcp"}, Help: "The total count of powerdns dns calls"})
	// powerdnsdnsudpcalltotal = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_powerdns_dns_call_total", ConstLabels: map[string]string{"state": "successful", "protocol": "udp"}, Help: "The total count of powerdns dns calls"}).
	// powerdnsdnsudperrtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_powerdns_dns_call_total", ConstLabels: map[string]string{"state": "failed", "protocol": "udp"}, Help: "The total count of powerdns dns calls"}).
	powerdnssoaquery     = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_powerdns_dns_query_total", ConstLabels: map[string]string{"state": "successful", "qtype": "soa"}, Help: "The total count of powerdns dns queries"})
	powerdnssoaqueryerr  = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_powerdns_dns_query_total", ConstLabels: map[string]string{"state": "failed", "qtype": "soa"}, Help: "The total count of powerdns dns queries"})
	powerdnsaxfrquery    = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_powerdns_dns_query_total", ConstLabels: map[string]string{"state": "successful", "qtype": "axfr"}, Help: "The total count of powerdns dns queries"})
	powerdnsaxfrqueryerr = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_powerdns_dns_query_total", ConstLabels: map[string]string{"state": "failed", "qtype": "axfr"}, Help: "The total count of powerdns dns queries"})
	responsetotal        = promauto.NewCounter(prometheus.CounterOpts{Name: "zoneprovider_response_total", Help: "The total count of responses"})
)

func ProcessRequest(msg *nats.Msg, conf *config.ServiceConfiguration) {
	zonename, wantsoa, parseerr := getZoneNameFrom(msg)
	if parseerr != nil {
		logger.ErrorErrLog(parseerr)

		return
	}

	logger.DebugLog("[Process Zone Data Request] for zone: " + zonename)

	zonedata, requesterr := getZoneDataFromPowerDNS(zonename, conf, wantsoa)
	if requesterr != nil {
		logger.ErrorErrLog(requesterr)

		return
	}

	dnssec, presigned, detailrequesterr := getZoneMetadataFromPowerDNS(zonename, conf)
	if detailrequesterr != nil {
		logger.ErrorErrLog(detailrequesterr)

		return
	}

	payload, marshalerr := createPayload(zonename, zonedata, dnssec, presigned)
	if marshalerr != nil {
		logger.ErrorErrLog(marshalerr)

		return
	}

	replyerr := replyOnRequest(msg, payload)
	if replyerr != nil {
		logger.ErrorErrLog(replyerr)

		return
	}
}

func getZoneNameFrom(msg *nats.Msg) (string, bool, error) {
	var zonerequest modelevent.ZoneDataRequestEvent

	incomingdata := msg.Data

	parseerr := json.Unmarshal(incomingdata, &zonerequest)
	if parseerr != nil {
		return "", false, parseerr
	}

	return zonerequest.Zone, zonerequest.OnlySOA, nil
}

func getZoneDataFromPowerDNS(zoneid string, conf *config.ServiceConfiguration, wantsoa bool) (zonefile string, err error) {
	readtimeout := conf.PowerDNSAXFRTimeout
	axfraddress := conf.PowerDNSAXFRAddress

	if wantsoa {
		zonefile, err = powerdnsutils.GetSOARecord(zoneid, axfraddress, "tcp")
		if err != nil {
			powerdnsdnstcperrtotal.Inc()
			powerdnssoaqueryerr.Inc()

			return
		}

		powerdnsdnstcpcalltotal.Inc()
		powerdnssoaquery.Inc()

		return
	}

	zonefile, err = powerdnsutils.GetZoneFilePerAXFR(axfraddress, zoneid, readtimeout)
	if err != nil {
		powerdnsdnstcperrtotal.Inc()
		powerdnsaxfrqueryerr.Inc()

		return
	}

	powerdnsdnstcpcalltotal.Inc()
	powerdnsaxfrquery.Inc()

	/* Use always AXFR for quering zones instead of api. Because LUA Records are not recognized by libs and pdns can interpret format
	if strings.Contains(zonefile, "CLASS1") {
		//CLASS1 is used by miek/dns due to not recognize LUA records
		zonefile, err = getZoneDataFromPowerDNSAPI(zoneid, conf)
		if err != nil {
			powerdnsapierrtotal.Inc()
			return
		}

		powerdnsapicalltotal.Inc()
	}
	*/

	return zonefile, nil
}

func getZoneDataFromPowerDNSAPI(zoneid string, conf *config.ServiceConfiguration) (string, error) { //nolint:unused,deadcode
	// does not work for dnssec signed zones
	// used for zones contain LUA records
	serverid := conf.PowerDNSServerID
	apitoken := conf.PowerDNSAPIToken
	zoneDataFileURL := conf.PowerDNSURL + "/api/v1/servers/" + serverid + "/zones/" + zoneid + "/export"

	return httputils.ExecutePowerDNSRequest(http.MethodGet, zoneDataFileURL, apitoken, nil)
}

func getZoneMetadataFromPowerDNS(zonename string, conf *config.ServiceConfiguration) (dnssec, presigned bool, err error) {
	zonedata, requesterr := getZoneDetailsPerPowerdnsAPI(conf, zonename)
	if requesterr != nil {
		powerdnsapierrtotal.Inc()

		return false, false, requesterr
	}

	presigned = false

	if zonedata.Dnssec {
		presigned = true
	}

	powerdnsapicalltotal.Inc()

	return zonedata.Dnssec, presigned, nil
}

func getZoneDetailsPerPowerdnsAPI(conf *config.ServiceConfiguration, zonename string) (modelpowerdns.Zone, error) {
	serverid := conf.PowerDNSServerID
	apitoken := conf.PowerDNSAPIToken
	zoneDataFileURL := conf.PowerDNSURL + "/api/v1/servers/" + serverid + "/zones/" + zonename

	detailresponse, httperr := httputils.ExecutePowerDNSRequest(http.MethodGet, zoneDataFileURL, apitoken, nil)
	if httperr != nil {
		return modelpowerdns.Zone{}, httperr
	}

	zonedata := modelpowerdns.Zone{}

	parseerr := json.Unmarshal([]byte(detailresponse), &zonedata)
	if parseerr != nil {
		return modelpowerdns.Zone{}, parseerr
	}

	return zonedata, nil
}

func createPayload(zonename, zonedata string, dnssec, presigned bool) ([]byte, error) {
	zonereply := modelevent.ZoneDataReplyEvent{
		Zone:             zonename,
		ZoneData:         zonedata,
		DnssecEnabled:    dnssec,
		PresignedRecords: presigned,
		RepliedAt:        time.Now(),
	}

	payload, parserr := json.Marshal(zonereply)
	if parserr != nil {
		return []byte{}, parserr
	}

	return payload, nil
}

func replyOnRequest(msg *nats.Msg, payload []byte) error {
	responsetotal.Inc()

	return msg.Respond(payload)
}

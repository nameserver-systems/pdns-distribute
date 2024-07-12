package powerdns

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/config"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelevent"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelpowerdns"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelzone"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/powerdnsutils"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var zonemapmutex sync.Map

const datarequesttimeout = 20 * time.Second

var (
	powerdnsapicreatecalltotal     = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_powerdns_api_call_total", ConstLabels: map[string]string{"state": "successful", "type": "create"}, Help: "The total count of powerdns api calls"})
	powerdnsapicreateerrtotal      = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_powerdns_api_call_total", ConstLabels: map[string]string{"state": "failed", "type": "create"}, Help: "The total count of powerdns api calls"})
	powerdnsapideletecalltotal     = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_powerdns_api_call_total", ConstLabels: map[string]string{"state": "successful", "type": "delete"}, Help: "The total count of powerdns api calls"})
	powerdnsapideleteerrtotal      = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_powerdns_api_call_total", ConstLabels: map[string]string{"state": "failed", "type": "delete"}, Help: "The total count of powerdns api calls"})
	powerdnsapigetzonescalltotal   = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_powerdns_api_call_total", ConstLabels: map[string]string{"state": "successful", "type": "get_zones"}, Help: "The total count of powerdns api calls"})
	powerdnsapigetzoneserrtotal    = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_powerdns_api_call_total", ConstLabels: map[string]string{"state": "failed", "type": "get_zones"}, Help: "The total count of powerdns api calls"})
	powerdnsapiclearcachecalltotal = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_powerdns_api_call_total", ConstLabels: map[string]string{"state": "successful", "type": "clear_cache"}, Help: "The total count of powerdns api calls"})

	zonedatarequestcalltotal = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_zone_data_request_total", ConstLabels: map[string]string{"state": "successful"}, Help: "The total count of zone data requests"})
	zonedatarequesterrtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_zone_data_request_total", ConstLabels: map[string]string{"state": "failed"}, Help: "The total count of zone data requests"})

	responsecalltotal = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_zone_state_response_total", ConstLabels: map[string]string{"state": "successful"}, Help: "The total count of zone state responses"})
	responseerrtotal  = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_zone_state_response_total", ConstLabels: map[string]string{"state": "failed"}, Help: "The total count of zone state responses"})

	workerlockcount = promauto.NewGauge(prometheus.GaugeOpts{Name: "secondarysyncer_worker_lock_count", Help: "The actual count of worker locks"})
)

func AddZone(msg *nats.Msg, conf *config.ServiceConfiguration) {
	pdnsconnection := modelpowerdns.PDNSconnectionobject{
		PowerDNSurl: conf.PowerDNSURL,
		ServerID:    conf.PowerDNSServerID,
		Apitoken:    conf.PowerDNSAPIToken,
	}

	zoneid, parseerr := getZoneIDFromAddEventMessage(msg)
	if parseerr != nil {
		logger.ErrorErrLog(parseerr)

		return
	}

	logger.DebugLog("[Add Zone] start triggered for zone: " + zoneid)

	zonemutex := startSynchronMutex(zoneid)
	defer stopSynchronMutex(zonemutex, zoneid)

	waitForPrimaryToStoreZoneAfterChange(conf)

	/* DISABLED DUE TO SLAVING ZONES | NECESSARY FOR USE WITH API OR ZONEFILES
	payload, marshalerr := PrepareZoneDataRequest(zoneid, false)
	if marshalerr != nil {
		logger.ErrorErrLog(marshalerr)
		return
	}

	zonedata, getdataerr := GetZoneData(conf, ms, payload)
	if getdataerr != nil {
		logger.ErrorErrLog(getdataerr)
		return
	}

	// necessary due specific cases, where zone some event arrives before the initial zone create event arrives
	// err will be suppressed, due to its not necessary in normal case
	if delerr := powerdnsutils.DeleteZone(pdnsconnection, zoneid); delerr != nil {
		powerdnsapideleteerrtotal.Inc()
	} else {
		powerdnsapideletecalltotal.Inc()
	} */

	zonedata := modelevent.ZoneDataReplyEvent{}

	createerr := createZone(zoneid, conf)
	if createerr != nil {
		logger.ErrorErrLog(createerr)

		return
	}

	if err := powerdnsutils.AXFRRetrieve(pdnsconnection, zoneid); err != nil {
		logger.ErrorErrLog(err)

		return
	}

	powerdnsutils.ClearCache(zoneid, pdnsconnection)
	powerdnsapiclearcachecalltotal.Inc()

	logger.DebugLog("[Add Zone] triggered for zone: " + zoneid + " with: " + zonedata.ZoneData)
}

func startSynchronMutex(zoneid string) interface{} {
	newzonemutex := new(sync.Mutex)
	zonemutex, _ := zonemapmutex.LoadOrStore(zoneid, newzonemutex)

	logger.DebugLog("[Synchronize Worker Lock] try to lock zone: " + zoneid)

	zonemutex.(*sync.Mutex).Lock()

	// // necessary for preventing race condition, but has more impact on cpu. alternative delete of zoneid in zonemapmutex can be ignored, but this has impact on ram => memory leak
	// zonemapmutex.Store(zoneid, zonemutex)

	workerlockcount.Inc()
	logger.DebugLog("[Synchronize Worker Lock] locked zone: " + zoneid)

	return zonemutex
}

func stopSynchronMutex(zonemutex interface{}, zoneid string) {
	zonemutex.(*sync.Mutex).Unlock()
	// zonemapmutex.Delete(zoneid)

	workerlockcount.Dec()
	logger.DebugLog("[Synchronize Worker Lock] unlocked zone: " + zoneid)
}

func createZone(zoneid string, conf *config.ServiceConfiguration) error { //nolint:wsl
	err := CreateZonePerAPI(zoneid, conf)
	if err != nil {
		return err
	}

	return nil
}

func CreateZonePerAPI(zoneid string, conf *config.ServiceConfiguration) error {
	pdnsconnection := modelpowerdns.PDNSconnectionobject{
		PowerDNSurl: conf.PowerDNSURL,
		ServerID:    conf.PowerDNSServerID,
		Apitoken:    conf.PowerDNSAPIToken,
	}

	payload, preparationerr := prepareCreateZoneRequest(zoneid)
	if preparationerr != nil {
		return preparationerr
	}

	logger.DebugLog("Create Zone payload" + string(payload))

	executeerr := powerdnsutils.CreateZone(pdnsconnection, payload)
	logger.DebugErrLog(executeerr)
	if executeerr != nil {
		powerdnsapicreateerrtotal.Inc()

		return executeerr
	}

	powerdnsapicreatecalltotal.Inc()

	return nil
}

func prepareCreateZoneRequest(zoneid string) ([]byte, error) {
	zonecreation := modelpowerdns.Zone{
		ID:          zoneid,
		Name:        zoneid,
		Kind:        "Slave",
		Masters:     []string{"127.0.0.1:20102"},
		Nameservers: make([]string, 0),
		//	Zone:        zonedata.ZoneData, // DISABLED DUE TO SLAVING ZONES
		SoaEdit:    "NONE",
		SoaEditAPI: "OFF",
		Presigned:  true,
	}

	storepayload, marshalerr := json.Marshal(zonecreation)
	if marshalerr != nil {
		return nil, marshalerr
	}

	return storepayload, nil
}

func GetZoneData(conf *config.ServiceConfiguration, ms *microservice.Microservice, payload []byte) (modelevent.ZoneDataReplyEvent, error) {
	zonedataresponse, err := sendZoneDataRequestAndWaitForResponse(conf, ms, payload)
	if err != nil {
		zonedatarequesterrtotal.Inc()

		return modelevent.ZoneDataReplyEvent{}, err
	}

	zonedatarequestcalltotal.Inc()

	zonedata := modelevent.ZoneDataReplyEvent{}

	responseunmarshalerr := json.Unmarshal(zonedataresponse, &zonedata)
	if responseunmarshalerr != nil {
		return modelevent.ZoneDataReplyEvent{}, responseunmarshalerr
	}

	return zonedata, nil
}

func sendZoneDataRequestAndWaitForResponse(conf *config.ServiceConfiguration, ms *microservice.Microservice, payload []byte) ([]byte, error) {
	sendertopic := conf.SenderTopic

	answer, requesterr := ms.MessageBroker.PublishRequestAndWait(sendertopic, payload, datarequesttimeout)
	if requesterr != nil {
		return nil, requesterr
	}

	zonedataresponse := answer.Data

	return zonedataresponse, nil
}

func PrepareZoneDataRequest(zoneid string, onlysoa bool) ([]byte, error) {
	zonerequest := modelevent.ZoneDataRequestEvent{
		Zone:        zoneid,
		OnlySOA:     onlysoa,
		RequestedAt: time.Now(),
	}

	payload, marshalerr := json.Marshal(zonerequest)
	if marshalerr != nil {
		return nil, marshalerr
	}

	return payload, nil
}

func waitForPrimaryToStoreZoneAfterChange(conf *config.ServiceConfiguration) {
	time.Sleep(conf.EventDelay)
}

func getZoneIDFromAddEventMessage(msg *nats.Msg) (string, error) {
	incomingmdata := msg.Data
	addevent := modelevent.ZoneAddEvent{}

	unmarshalerr := json.Unmarshal(incomingmdata, &addevent)
	if unmarshalerr != nil {
		logger.ErrorErrLog(unmarshalerr)

		return "", unmarshalerr
	}

	return addevent.Zone, nil
}

func getZoneIDFromChangeEventMessage(msg *nats.Msg) (string, error) {
	incomingmdata := msg.Data
	changeevent := modelevent.ZoneChangeEvent{}

	unmarshalerr := json.Unmarshal(incomingmdata, &changeevent)
	if unmarshalerr != nil {
		logger.ErrorErrLog(unmarshalerr)

		return "", unmarshalerr
	}

	return changeevent.Zone, nil
}

func getZoneIDFromDeleteEventMessage(msg *nats.Msg) (string, error) {
	incomingmdata := msg.Data
	deleteevent := modelevent.ZoneDeleteEvent{}

	unmarshalerr := json.Unmarshal(incomingmdata, &deleteevent)
	if unmarshalerr != nil {
		logger.ErrorErrLog(unmarshalerr)

		return "", unmarshalerr
	}

	return deleteevent.Zone, nil
}

func ChangeZone(msg *nats.Msg, conf *config.ServiceConfiguration) { //nolint:funlen
	pdnsconnection := modelpowerdns.PDNSconnectionobject{
		PowerDNSurl: conf.PowerDNSURL,
		ServerID:    conf.PowerDNSServerID,
		Apitoken:    conf.PowerDNSAPIToken,
	}

	zoneid, parseerr := getZoneIDFromChangeEventMessage(msg)
	if parseerr != nil {
		logger.ErrorErrLog(parseerr)

		return
	}

	logger.DebugLog("[Change Zone] start triggered for zone: " + zoneid)

	zonemutex := startSynchronMutex(zoneid)
	defer stopSynchronMutex(zonemutex, zoneid)

	waitForPrimaryToStoreZoneAfterChange(conf)

	/* DISABLED DUE TO SLAVING ZONES | NECESSARY FOR USE WITH API OR ZONEFILES

	payload, marshalerr := PrepareZoneDataRequest(zoneid)
	if marshalerr != nil {
		logger.ErrorErrLog(marshalerr)
		return
	}

	zonedata, getdataerr := GetZoneData(conf, ms, payload)
	if getdataerr != nil {
		logger.ErrorErrLog(getdataerr)
		return
	}

	delerr := powerdnsutils.DeleteZone(pdnsconnection, zoneid)
	if delerr != nil {
		//No return after logging, because its allowed to fail here for zone creation
		logger.ErrorErrLog(delerr)
		powerdnsapideleteerrtotal.Inc()
	} else {
		powerdnsapideletecalltotal.Inc()
	}
	*/

	zoneExist, err := powerdnsutils.DoesZoneExist(pdnsconnection, zoneid)
	if err != nil {
		logger.ErrorErrLog(err)

		return
	}

	if !zoneExist {
		createerr := createZone(zoneid, conf)
		if createerr != nil {
			logger.ErrorErrLog(createerr)

			return
		}
	}

	if err := powerdnsutils.AXFRRetrieve(pdnsconnection, zoneid); err != nil {
		logger.ErrorErrLog(err)
	}

	powerdnsutils.ClearCache(zoneid, pdnsconnection)
	powerdnsapiclearcachecalltotal.Inc()

	logger.DebugLog("[Change Zone] triggered for zone: " + zoneid + " with: AXFR") // zonedata.ZoneData)
}

func DeleteZone(msg *nats.Msg, conf *config.ServiceConfiguration) {
	pdnsconnection := modelpowerdns.PDNSconnectionobject{
		PowerDNSurl: conf.PowerDNSURL,
		ServerID:    conf.PowerDNSServerID,
		Apitoken:    conf.PowerDNSAPIToken,
	}

	zoneid, parseerr := getZoneIDFromDeleteEventMessage(msg)
	if parseerr != nil {
		logger.ErrorErrLog(parseerr)

		return
	}

	logger.DebugLog("[Delete Zone] start triggered for zone: " + zoneid)

	zonemutex := startSynchronMutex(zoneid)
	defer stopSynchronMutex(zonemutex, zoneid)

	delerr := powerdnsutils.DeleteZone(pdnsconnection, zoneid)
	if delerr != nil {
		powerdnsapideleteerrtotal.Inc()
		logger.ErrorErrLog(delerr)

		return
	}

	powerdnsapideletecalltotal.Inc()

	powerdnsutils.ClearCache(zoneid, pdnsconnection)
	powerdnsapiclearcachecalltotal.Inc()

	logger.DebugLog("[Delete Zone] triggered for zone: " + zoneid)
}

func RespondSecondaryZoneState(msg *nats.Msg, conf *config.ServiceConfiguration) {
	var actualzones []modelpowerdns.Zone

	var request modelzone.Zonestaterequestevent

	pdnsconnection := modelpowerdns.PDNSconnectionobject{
		PowerDNSurl: conf.PowerDNSURL,
		ServerID:    conf.PowerDNSServerID,
		Apitoken:    conf.PowerDNSAPIToken,
	}

	incomingrequest := msg.Data

	unmarshalerr := json.Unmarshal(incomingrequest, &request)
	if unmarshalerr != nil {
		logger.ErrorErrLog(unmarshalerr)
	}

	zoneid := request.Zone

	response, err := powerdnsutils.GetZoneList(pdnsconnection, zoneid, false)
	if err != nil {
		logger.ErrorErrLog(err)
		powerdnsapigetzoneserrtotal.Inc()
	} else {
		powerdnsapigetzonescalltotal.Inc()
	}

	parseerr := json.Unmarshal([]byte(response), &actualzones)
	if parseerr != nil {
		logger.ErrorErrLog(parseerr)
	}

	zonestatemap := powerdnsutils.TransferPowerDNSZonesIntoZoneStateMap(actualzones)

	transferobject := modelzone.Zonestatemaptransferobject{
		Statemap:  zonestatemap,
		CreatedAT: time.Now(),
	}

	responsedata, marshalerr := json.Marshal(transferobject)
	if marshalerr != nil {
		logger.ErrorErrLog(marshalerr)
	}

	responderr := msg.Respond(responsedata)
	if responderr != nil {
		logger.ErrorErrLog(responderr)
		responseerrtotal.Inc()
	} else {
		responsecalltotal.Inc()
	}
}

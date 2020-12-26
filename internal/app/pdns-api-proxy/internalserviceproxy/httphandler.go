package internalserviceproxy

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/eventutils"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/httputils"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

var (
	powerdnsapirequeststotal = promauto.NewCounter(prometheus.CounterOpts{Name: "apiproxy_powerdns_api_requests_total", Help: "The total count of powerdns api requests"})
	changepublishedtotal     = promauto.NewCounter(prometheus.CounterOpts{Name: "apiproxy_published_zone_events_total", ConstLabels: map[string]string{"event_type": "change"}, Help: "The total count of published zone events"})
	deletepublishedtotal     = promauto.NewCounter(prometheus.CounterOpts{Name: "apiproxy_published_zone_events_total", ConstLabels: map[string]string{"event_type": "delete"}, Help: "The total count of published zone events"})
	createpublishedtotal     = promauto.NewCounter(prometheus.CounterOpts{Name: "apiproxy_published_zone_events_total", ConstLabels: map[string]string{"event_type": "create"}, Help: "The total count of published zone events"})
)

func createZoneHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := httputils.HTTPResponseStatusWriter{ResponseWriter: w}

		zoneid, reuqestparseerr := httputils.GetZoneIDFromRequest(r)

		f.ServeHTTP(&response, r)

		if reuqestparseerr != nil {
			logger.ErrorErrLog(reuqestparseerr)

			return
		}

		if httputils.IsStatusCodeSuccesful(response.Status) {
			eventutils.PublishCreateZoneEvent(ms, serviceconfig.AddEventTopic, zoneid)
			createpublishedtotal.Inc()
		}
	}
}

func changeZoneHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := httputils.HTTPResponseStatusWriter{ResponseWriter: w}

		zoneid, reuqestparseerr := httputils.GetZoneIDFromRequest(r)

		f.ServeHTTP(&response, r)

		if reuqestparseerr != nil {
			logger.ErrorErrLog(reuqestparseerr)

			return
		}

		if httputils.IsStatusCodeSuccesful(response.Status) {
			eventutils.PublishChangeZoneEvent(ms, serviceconfig.ChangeEventTopic, zoneid)
			changepublishedtotal.Inc()
		}
	}
}

func deleteZoneHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := httputils.HTTPResponseStatusWriter{ResponseWriter: w}

		zoneid, reuqestparseerr := httputils.GetZoneIDFromRequest(r)

		f.ServeHTTP(&response, r)

		if reuqestparseerr != nil {
			logger.ErrorErrLog(reuqestparseerr)

			return
		}

		if httputils.IsStatusCodeSuccesful(response.Status) {
			eventutils.PublishDeleteZoneEvent(ms, serviceconfig.DeleteEventTopic, zoneid)
			deletepublishedtotal.Inc()
		}
	}
}

func defaultProxy(respond http.ResponseWriter, incomingrequest *http.Request) {
	var proxyclient http.Client

	powerdnsapirequeststotal.Inc()

	incomingbodybytes, readerr := ioutil.ReadAll(incomingrequest.Body)
	if readerr != nil {
		logger.ErrorErrLog(readerr)
	}

	incomingbodyreader := bytes.NewReader(incomingbodybytes)

	proxyurl := serviceconfig.PowerDNSURL + incomingrequest.URL.Path
	if incomingrequest.URL.RawQuery != "" {
		proxyurl += "?" + incomingrequest.URL.RawQuery
	}

	proxyrequest, preparerequesterr := http.NewRequestWithContext(context.Background(), incomingrequest.Method, proxyurl, incomingbodyreader)
	if preparerequesterr != nil {
		logger.ErrorErrLog(preparerequesterr)
	}

	httputils.CopyHTTPHeader(incomingrequest.Header, proxyrequest.Header)

	proxyresponse, proxyexecuteerr := proxyclient.Do(proxyrequest) //nolint:bodyclose
	if proxyexecuteerr != nil {
		logger.ErrorErrLog(proxyexecuteerr)
		respond.WriteHeader(http.StatusBadGateway)

		return
	}

	defer httputils.CloseResponseBody(proxyresponse)

	proxyresponsebody, proxyresponsereaderr := ioutil.ReadAll(proxyresponse.Body)
	if proxyresponsereaderr != nil {
		logger.ErrorErrLog(proxyresponsereaderr)
	}

	httputils.CopyHTTPHeader(proxyresponse.Header, respond.Header())

	respond.WriteHeader(proxyresponse.StatusCode)

	_, readerr = respond.Write(proxyresponsebody)
	if readerr != nil {
		logger.ErrorErrLog(readerr)
	}
}

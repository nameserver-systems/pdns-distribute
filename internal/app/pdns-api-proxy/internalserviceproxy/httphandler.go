package internalserviceproxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/nameserver-systems/pdns-distribute/internal/pkg/eventutils"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/httputils"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

	incomingbodybytes, readerr := io.ReadAll(incomingrequest.Body)
	if readerr != nil {
		logger.ErrorErrLog(readerr)
	}

	incomingbodyreader := bytes.NewReader(incomingbodybytes)

	proxyurl := serviceconfig.PowerDNSURL + incomingrequest.URL.Path
	if incomingrequest.URL.RawQuery != "" {
		proxyurl += "?" + incomingrequest.URL.RawQuery
	}

	proxyrequest, preparerequesterr := http.NewRequestWithContext(incomingrequest.Context(), incomingrequest.Method, proxyurl, incomingbodyreader)
	if preparerequesterr != nil {
		logger.ErrorErrLog(preparerequesterr)
	}

	proxyrequest.Header = incomingrequest.Header.Clone()

	proxyresponse, proxyexecuteerr := proxyclient.Do(proxyrequest) //nolint:bodyclose
	if proxyexecuteerr != nil {
		logger.ErrorErrLog(proxyexecuteerr)
		respond.WriteHeader(http.StatusBadGateway)

		return
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.ErrorErrLog(err)
		}
	}(proxyresponse.Body)

	proxyresponsebody, proxyresponsereaderr := io.ReadAll(proxyresponse.Body)
	if proxyresponsereaderr != nil {
		logger.ErrorErrLog(proxyresponsereaderr)
	}

	for key, values := range proxyresponse.Header {
		for _, value := range values {
			respond.Header().Add(key, value)
		}
	}

	respond.WriteHeader(proxyresponse.StatusCode)

	_, readerr = respond.Write(proxyresponsebody)
	if readerr != nil {
		logger.ErrorErrLog(readerr)
	}
}

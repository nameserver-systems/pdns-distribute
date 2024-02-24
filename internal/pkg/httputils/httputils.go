package httputils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	neturl "net/url"
	"strconv"

	"github.com/miekg/dns"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

func GetZoneIDFromRequest(r *http.Request) (string, error) {
	var err error

	zoneid := getZoneIDFromRequestPath(r)

	if len(zoneid) == 0 {
		zoneid, err = getZoneIDFromRequestBody(r)
		if err != nil {
			return "", err
		}
	}

	zoneid = dns.Fqdn(zoneid)

	return zoneid, nil
}

func getZoneIDFromRequestBody(r *http.Request) (zoneID string, err error) {
	var result map[string]interface{}

	incomingbodybytes, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	if err = r.Body.Close(); err != nil {
		return
	}

	// necessary due to io.ReadAll clears after read the body and body is used for proxy
	r.Body = io.NopCloser(bytes.NewBuffer(incomingbodybytes))

	if err = json.Unmarshal(incomingbodybytes, &result); err != nil {
		return
	}

	if rawzoneid, okid := result["id"]; okid {
		zoneID = rawzoneid.(string)
	} else if rawzonename, okname := result["name"]; okname {
		zoneID = rawzonename.(string)
	}

	if len(zoneID) == 0 {
		err = errZoneIDNotFound
	}

	return
}

func getZoneIDFromRequestPath(r *http.Request) string {
	return r.PathValue("zone_id")
}

func GetHostAndPortFromURL(url string) (string, error) {
	urlobject, parseerr := neturl.Parse(url)
	if parseerr != nil {
		urlobject = nil
	}

	address := urlobject.Host
	if address == "" {
		return "", errEmptyHostAddress
	}

	return address, parseerr
}

func GetHostnameFromURL(url string) (string, error) {
	urlobject, parseerr := neturl.Parse(url)
	if parseerr != nil {
		urlobject = nil
	}

	hostname := urlobject.Hostname()
	if hostname == "" {
		return "", errEmptyHostname
	}

	return hostname, parseerr
}

func ExecutePowerDNSRequest(method, url, apitoken string, body io.Reader) (string, error) {
	var client http.Client

	request, preparerequesterr := http.NewRequestWithContext(context.Background(), method, url, body)
	if preparerequesterr != nil {
		return "", preparerequesterr
	}

	request.Header.Add("X-API-Key", apitoken)

	response, executeerr := client.Do(request)
	if executeerr != nil {
		return "", executeerr
	}

	defer response.Body.Close()

	responsebody, responsereaderr := io.ReadAll(response.Body)
	if responsereaderr != nil {
		return "", responsereaderr
	}

	if !IsStatusCodeSuccesful(response.StatusCode) {
		logger.ErrorLog("invalid pdns answer with http status" + strconv.Itoa(response.StatusCode) + " with message: " + string(responsebody))

		return "", errInvalidPDNSAPIAnswer
	}

	return string(responsebody), nil
}

func IsStatusCodeSuccesful(statuscode int) bool {
	return statuscode >= 200 && statuscode < 300
}

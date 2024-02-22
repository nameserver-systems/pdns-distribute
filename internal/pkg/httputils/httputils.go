package httputils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

func CopyHTTPHeader(src, dest http.Header) {
	for key, values := range src {
		for _, value := range values {
			dest.Add(key, value)
		}
	}
}

func GetZoneIDFromRequest(r *http.Request) (string, error) {
	var err error

	zoneid := getZoneIDFromRequestPath(r)

	if hasNotFoundZoneID(zoneid) {
		zoneid, err = getZoneIDFromRequestBody(r, zoneid)
		if err != nil {
			return "", err
		}
	}

	zoneid = ensureTrailingDot(zoneid)

	return zoneid, nil
}

func ensureTrailingDot(zoneid string) string {
	if !strings.HasSuffix(zoneid, ".") {
		zoneid += "."
	}

	return zoneid
}

func getZoneIDFromRequestBody(r *http.Request, zoneid string) (string, error) {
	var result map[string]interface{}

	incomingbodybytes, readerr := io.ReadAll(r.Body)
	if readerr != nil {
		return "", readerr
	}

	closerr := r.Body.Close()
	if closerr != nil {
		return "", closerr
	}

	// necessary due to io.ReadAll clears after read the body and body is used for proxy
	r.Body = io.NopCloser(bytes.NewBuffer(incomingbodybytes))

	unmarshalerr := json.Unmarshal(incomingbodybytes, &result)
	if unmarshalerr != nil {
		return "", unmarshalerr
	}

	if rawzoneid, okid := result["id"]; okid {
		zoneid = rawzoneid.(string)
	} else if rawzonename, okname := result["name"]; okname {
		zoneid = rawzonename.(string)
	}

	if hasNotFoundZoneID(zoneid) {
		return "", errZoneIDNotFound
	}

	return zoneid, nil
}

func getZoneIDFromRequestPath(r *http.Request) string {
	vars := mux.Vars(r)

	zoneid := vars["zone_id"]

	return zoneid
}

func hasNotFoundZoneID(zoneid string) bool {
	return zoneid == ""
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

func CloseResponseBody(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		logger.ErrorErrLog(err)
	}
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

	defer CloseResponseBody(response)

	responsebody, responsereaderr := io.ReadAll(response.Body)
	if responsereaderr != nil {
		return "", responsereaderr
	}

	if 200 > response.StatusCode || response.StatusCode >= 300 {
		logger.ErrorLog("invalid pdns answer with http status" + strconv.Itoa(response.StatusCode) + " with message: " + string(responsebody))

		return "", errInvalidPDNSAPIAnswer
	}

	return string(responsebody), nil
}

func IsStatusCodeSuccesful(statuscode int) bool {
	return statuscode >= 200 && statuscode < 300
}

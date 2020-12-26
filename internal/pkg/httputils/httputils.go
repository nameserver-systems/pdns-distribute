package httputils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
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

	incomingbodybytes, readerr := ioutil.ReadAll(r.Body)
	if readerr != nil {
		return "", readerr
	}

	closerr := r.Body.Close()
	if closerr != nil {
		return "", closerr
	}

	// necessary due to ioutil.readall clears after read the body and body is used for proxy
	r.Body = ioutil.NopCloser(bytes.NewBuffer(incomingbodybytes))

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
		return "", errors.New("can't find ZoneID")
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
		return "", errors.New("parse error: empty extracted address")
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
		return "", errors.New("parse error: empty extracted hostname")
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

	response, executeerr := client.Do(request) //nolint:bodyclose
	if executeerr != nil {
		return "", executeerr
	}

	defer CloseResponseBody(response)

	responsebody, responsereaderr := ioutil.ReadAll(response.Body)
	if responsereaderr != nil {
		return "", responsereaderr
	}

	if 200 > response.StatusCode || response.StatusCode >= 300 {
		return "", errors.New("couldn't get valid answer from powerdns api, http status is " +
			strconv.Itoa(response.StatusCode) + " with message: " + string(responsebody))
	}

	return string(responsebody), nil
}

func IsStatusCodeSuccesful(statuscode int) bool {
	return statuscode >= 200 && statuscode < 300
}

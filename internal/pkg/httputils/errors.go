package httputils

import "errors"

var (
	errZoneIDNotFound       = errors.New("can't find ZoneID")
	errEmptyHostAddress     = errors.New("parse error: empty extracted address")
	errEmptyHostname        = errors.New("parse error: empty extracted hostname")
	errInvalidPDNSAPIAnswer = errors.New("couldn't get valid answer from powerdns api")
)

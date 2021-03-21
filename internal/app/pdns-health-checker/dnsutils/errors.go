package dnsutils

import "errors"

var (
	errTooMuchZones   = errors.New("got more zones than expected")
	errZoneIDMismatch = errors.New("zoneid mismatch, can't get serial")
)

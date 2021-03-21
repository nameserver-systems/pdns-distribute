package utils

import "errors"

var (
	errUnexpectedZoneCount = errors.New("got unexpected zone count")
	errZoneIDNotInStateMap = errors.New("zoneid not found in statemap")
)

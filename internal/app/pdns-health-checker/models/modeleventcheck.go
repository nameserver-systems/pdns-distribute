package models

const (
	AddZone = iota
	ChangeZone
	DeleteZone
)

type EventCheckObject struct {
	Eventtype     int
	Zoneid        string
	Primaryserial int32
}

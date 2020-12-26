package modelzone

import (
	"time"
)

// Zonestaterequestevent requests a ZoneState. If zone is empty it should return a map for all zones.
type Zonestaterequestevent struct {
	Zone        string    `json:"zone,omitempty"`
	RequestedAt time.Time `json:"requested_at,omitempty"`
}

type Zonestatemaptransferobject struct {
	Statemap  Zonestatemap `json:"statemap"`
	CreatedAT time.Time    `json:"created_at,omitempty"`
}

// Zonestatemap key = zoneid, value zone serial.
type Zonestatemap map[string]int32

func (primaryzones Zonestatemap) Diff(secondaryzones Zonestatemap) (addedzones, deletedzones, changedzones Zonestatemap) {
	addedzones = make(Zonestatemap)
	deletedzones = make(Zonestatemap)
	changedzones = make(Zonestatemap)

	for zoneid, serial := range primaryzones {
		secondaryserial, exist := secondaryzones[zoneid]
		if exist {
			if serial != secondaryserial {
				changedzones[zoneid] = serial

				delete(secondaryzones, zoneid)
			}
		} else {
			addedzones[zoneid] = serial
			delete(secondaryzones, zoneid)
		}
	}

	for zoneid, secondaryserial := range secondaryzones {
		_, exist := primaryzones[zoneid]
		if !exist {
			deletedzones[zoneid] = secondaryserial

			delete(secondaryzones, zoneid)
		}
	}

	return
}

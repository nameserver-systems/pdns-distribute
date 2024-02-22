package models

import (
	"sync"
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelzone"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/messaging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	zonestatelockcount      = promauto.NewGauge(prometheus.GaugeOpts{Name: "healthchecker_zone_state_lock_count", Help: "The actual count of zone state locks"})
	secondarystatelockcount = promauto.NewGauge(prometheus.GaugeOpts{Name: "healthchecker_secondary_state_lock_count", Help: "The actual count of secondary state locks"})
)

type State struct {
	zoneMutex      sync.Mutex
	expectedZones  modelzone.Zonestatemap
	zonesChangedAt time.Time

	secondaryMutex       sync.Mutex
	activeSecondaries    []messaging.ResolvedService
	secondariesChangedAt time.Time
}

func GenerateStateObject() *State {
	return &State{}
}

func (s *State) SetExpectedZoneMap(expectedzones modelzone.Zonestatemap) {
	s.zoneMutex.Lock()
	zonestatelockcount.Inc()

	defer func() {
		s.zoneMutex.Unlock()
		zonestatelockcount.Dec()
	}()

	s.expectedZones = expectedzones
	s.zonesChangedAt = time.Now()
}

func (s *State) GetExpectedZoneMap() modelzone.Zonestatemap {
	s.zoneMutex.Lock()
	zonestatelockcount.Inc()

	defer func() {
		s.zoneMutex.Unlock()
		zonestatelockcount.Dec()
	}()

	return s.expectedZones
}

func (s *State) GetExpectedZoneChangeTime() time.Time {
	s.zoneMutex.Lock()
	zonestatelockcount.Inc()

	defer func() {
		s.zoneMutex.Unlock()
		zonestatelockcount.Dec()
	}()

	return s.zonesChangedAt
}

func (s *State) SetActiveSecondaries(secondaries []messaging.ResolvedService) {
	s.secondaryMutex.Lock()
	secondarystatelockcount.Inc()

	defer func() {
		s.secondaryMutex.Unlock()
		secondarystatelockcount.Dec()
	}()

	s.activeSecondaries = secondaries
	s.secondariesChangedAt = time.Now()
}

func (s *State) GetActiveSecondaries() []messaging.ResolvedService {
	s.secondaryMutex.Lock()
	secondarystatelockcount.Inc()

	defer func() {
		s.secondaryMutex.Unlock()
		secondarystatelockcount.Dec()
	}()

	return s.activeSecondaries
}

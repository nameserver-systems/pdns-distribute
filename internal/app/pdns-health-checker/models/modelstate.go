package models

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/modelzone"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/servicediscovery"
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
	activeSecondaries    []servicediscovery.ResolvedService
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

func (s *State) SetActiveSecondaries(secondaries []servicediscovery.ResolvedService) {
	s.secondaryMutex.Lock()
	secondarystatelockcount.Inc()

	defer func() {
		s.secondaryMutex.Unlock()
		secondarystatelockcount.Dec()
	}()

	s.activeSecondaries = secondaries
	s.secondariesChangedAt = time.Now()
}

func (s *State) GetActiveSecondaries() []servicediscovery.ResolvedService {
	s.secondaryMutex.Lock()
	secondarystatelockcount.Inc()

	defer func() {
		s.secondaryMutex.Unlock()
		secondarystatelockcount.Dec()
	}()

	return s.activeSecondaries
}

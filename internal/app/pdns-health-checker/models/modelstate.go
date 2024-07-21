package models

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/nameserver-systems/pdns-distribute/internal/pkg/modelzone"
	msframe "github.com/nameserver-systems/pdns-distribute/pkg/microservice"
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

	secondaryMutex sync.Mutex
	secondaries    []HashedSecondary
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

type HashedSecondary struct {
	ID               string
	SecondaryAddress string
}

func (s *State) SetActiveSecondaries(microservice *msframe.Microservice) {
	s.secondaryMutex.Lock()
	secondarystatelockcount.Inc()

	defer func() {
		s.secondaryMutex.Unlock()
		secondarystatelockcount.Dec()
	}()

	s.secondaries = make([]HashedSecondary, len(microservice.Secondaries))
	for i, secondary := range microservice.Secondaries {
		h := HashedSecondary{
			ID:               fmt.Sprintf("%x", sha256.Sum256([]byte(secondary))),
			SecondaryAddress: secondary,
		}
		s.secondaries[i] = h
	}
}

func (s *State) GetActiveSecondaries() []HashedSecondary {
	s.secondaryMutex.Lock()
	secondarystatelockcount.Inc()

	defer func() {
		s.secondaryMutex.Unlock()
		secondarystatelockcount.Dec()
	}()

	return s.secondaries
}

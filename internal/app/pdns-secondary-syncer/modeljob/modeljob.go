package modeljob

import (
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/config"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	AddZone = iota
	ChangeZone
	DeleteZone
)

type PowerDNSAPIJob struct {
	Jobtype int
	Msg     jetstream.Msg
	Ms      *microservice.Microservice
	Conf    *config.ServiceConfiguration
}

package modeljob

import (
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/config"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nats-io/nats.go"
)

const (
	AddZone = iota
	ChangeZone
	DeleteZone
)

type PowerDNSAPIJob struct {
	Jobtype int
	Msg     *nats.Msg
	Ms      *microservice.Microservice
	Conf    *config.ServiceConfiguration
}

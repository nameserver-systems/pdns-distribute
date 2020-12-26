package modeljob

import (
	"github.com/nats-io/nats.go"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/config"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice"
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

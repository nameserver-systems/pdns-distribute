package externalserviceproxy

import (
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-zone-provider/config"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-zone-provider/messaging"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
)

func StartExternalServiceProxy(ms *microservice.Microservice) error {
	serviceconfig := config.GetConfiguration(ms)

	messaging.SubscribeToIncomingMessages(ms, serviceconfig)

	return nil
}

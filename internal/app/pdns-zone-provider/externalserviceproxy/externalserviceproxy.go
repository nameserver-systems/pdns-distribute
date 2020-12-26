package externalserviceproxy

import (
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-zone-provider/config"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-zone-provider/messaging"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice"
)

func StartExternalServiceProxy(ms *microservice.Microservice) error {
	serviceconfig := config.GetConfiguration(ms)

	messaging.SubscribeToIncomingMessages(ms, serviceconfig)

	return nil
}

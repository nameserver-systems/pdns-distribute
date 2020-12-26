package servicediscovery

import (
	"net/url"
	"time"

	consul "github.com/hashicorp/consul/api"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/utils"
)

// ServiceRegistration used for register service in service discovery.
type ServiceRegistration struct {
	ServiceDiscoveryURL                 string
	ServiceDiscoveryUsername            string
	ServiceDiscoveryPassword            string
	ServiceDiscoveryHealthPingIntervall time.Duration

	MicroserviceID       string
	MicroserviceName     string
	MicroserviceTags     []string
	MicroserviceMetadata map[string]string
	MicroserviceURL      string
}

func (sr *ServiceRegistration) isURLSet() bool {
	return sr.ServiceDiscoveryURL != ""
}

func (sr *ServiceRegistration) generateServiceDiscoveryRegistration() (consul.AgentServiceRegistration, error) {
	microserviceurl := sr.MicroserviceURL

	msurl, parseerr := url.Parse(microserviceurl)
	if parseerr != nil {
		return consul.AgentServiceRegistration{}, parseerr
	}

	microserviceport, converterr := utils.ConvertStringToInt(msurl.Port())
	if converterr != nil {
		return consul.AgentServiceRegistration{}, converterr
	}

	return consul.AgentServiceRegistration{
		ID:      sr.MicroserviceID,
		Name:    sr.MicroserviceName,
		Tags:    sr.MicroserviceTags,
		Port:    microserviceport,
		Address: msurl.Hostname(),
		Meta:    sr.MicroserviceMetadata,
	}, nil
}

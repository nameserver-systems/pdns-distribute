package servicediscovery

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
	"github.com/shirou/gopsutil/load"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/utils"
)

// ServiceDiscovery struct which holds connection.
type ServiceDiscovery struct {
	serviceID string

	consulURL string

	consulBasicAuth *consul.HttpBasicAuth
	ConsulClient    *consul.Client

	serviceListenAddresses map[string]consul.ServiceAddress
	serviceHealthCheck     *consul.AgentServiceCheck
}

// StartServiceDiscoveryAndRegisterService starts discovery connection ang registers service.
func StartServiceDiscoveryAndRegisterService(sr *ServiceRegistration) (*ServiceDiscovery, error) {
	sd := &ServiceDiscovery{}

	err := sd.startServiceDiscoveryConnection(sr)
	if err != nil {
		return &ServiceDiscovery{}, err
	}

	return sd, nil
}

func (sd *ServiceDiscovery) startServiceDiscoveryConnection(sr *ServiceRegistration) error {
	if !sr.isURLSet() {
		return errNoURL
	}

	sd.insertServiceRegistration(sr)

	if sr.ServiceDiscoveryUsername != "" && sr.ServiceDiscoveryPassword != "" {
		sd.setServiceDiscoveryBasicAuthCredentials(sr.ServiceDiscoveryUsername, sr.ServiceDiscoveryPassword)
	}

	servicereg, err := sr.generateServiceDiscoveryRegistration()
	if err != nil {
		return err
	}

	addaddresserr := sd.addServiceListenAddress(sr.MicroserviceURL, "MainService")
	if addaddresserr != nil {
		return addaddresserr
	}

	healthreg := generateHealthCheckRegistration(sr)
	sd.setHealthCheck(healthreg)

	regerr := sd.registerService(servicereg)

	sd.startUpdateHealthCheckTTLHandler(sr)

	return regerr
}

func (sd *ServiceDiscovery) startUpdateHealthCheckTTLHandler(sr *ServiceRegistration) {
	healthcheckping := time.Tick(sr.ServiceDiscoveryHealthPingIntervall) //nolint:staticcheck
	checkid := sd.serviceHealthCheck.CheckID

	go func() {
		for range healthcheckping {
			status := consul.HealthPassing

			if isLoadTooHigh() {
				status = consul.HealthWarning
			}

			updateerr := sd.ConsulClient.Agent().UpdateTTL(checkid, "", status)
			if updateerr != nil {
				logger.ErrorErrLog(updateerr)
			}
		}
	}()
}

func isLoadTooHigh() bool {
	const highload = 10.0

	actualload, loaderr := load.Avg()
	if loaderr != nil {
		logger.ErrorErrLog(loaderr)
	}

	return actualload.Load1 > highload
}

func (sd *ServiceDiscovery) insertServiceRegistration(sr *ServiceRegistration) {
	sd.serviceID = sr.MicroserviceID
	sd.consulURL = sr.ServiceDiscoveryURL
}

func generateHealthCheckRegistration(sr *ServiceRegistration) *consul.AgentServiceCheck {
	return &consul.AgentServiceCheck{
		CheckID:                        uuid.NewString(),
		TTL:                            (sr.ServiceDiscoveryHealthPingIntervall * 3).String(),
		DeregisterCriticalServiceAfter: (sr.ServiceDiscoveryHealthPingIntervall * 20).String(),
	}
}

func (sd *ServiceDiscovery) addServiceListenAddress(serviceurl string, description string) error {
	if len(sd.serviceListenAddresses) == 0 {
		sd.serviceListenAddresses = make(map[string]consul.ServiceAddress)
	}

	svurl, err := url.Parse(serviceurl)
	if err != nil {
		return err
	}

	numericport, convrsionerr := utils.ConvertStringToInt(svurl.Port())
	if convrsionerr != nil {
		return convrsionerr
	}

	sd.serviceListenAddresses[description] = consul.ServiceAddress{
		Address: svurl.Hostname(),
		Port:    numericport,
	}

	return nil
}

func (sd *ServiceDiscovery) setHealthCheck(healthCheck *consul.AgentServiceCheck) {
	sd.serviceHealthCheck = healthCheck
}

func (sd *ServiceDiscovery) isListenAddressesSet() bool {
	return len(sd.serviceListenAddresses) != 0
}

func (sd *ServiceDiscovery) isClientInitiated() bool {
	return sd.ConsulClient != nil
}

func (sd *ServiceDiscovery) isHealthCheckSet() bool {
	return sd.serviceHealthCheck != nil
}

func (sd *ServiceDiscovery) registerService(serviceRegistration consul.AgentServiceRegistration) error {
	clienterr := sd.newClient()
	if clienterr != nil {
		return clienterr
	}

	if sd.isListenAddressesSet() {
		sd.insertListenAddressesInRegistration(&serviceRegistration)
	}

	if sd.isHealthCheckSet() {
		sd.insertHealthCheckInRegistration(&serviceRegistration)
	}

	registererr := sd.ConsulClient.Agent().ServiceRegister(&serviceRegistration)
	if registererr != nil {
		return registererr
	}

	return nil
}

// DeregisterService used with defer to deregister service if graceful shutdown.
func (sd *ServiceDiscovery) DeregisterService() error {
	if sd.ConsulClient != nil {
		err := sd.ConsulClient.Agent().ServiceDeregister(sd.serviceID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sd *ServiceDiscovery) insertListenAddressesInRegistration(serviceRegistration *consul.AgentServiceRegistration) {
	serviceRegistration.TaggedAddresses = sd.serviceListenAddresses
}

func (sd *ServiceDiscovery) insertHealthCheckInRegistration(serviceRegistration *consul.AgentServiceRegistration) {
	serviceRegistration.Check = sd.serviceHealthCheck
}

func (sd *ServiceDiscovery) setServiceDiscoveryBasicAuthCredentials(username, password string) {
	sd.consulBasicAuth = &consul.HttpBasicAuth{
		Username: username,
		Password: password,
	}
}

func (sd *ServiceDiscovery) newClient() error {
	if sd.isClientInitiated() {
		return nil
	}

	conf := consul.Config{
		Address: sd.consulURL,
	}

	if sd.consulBasicAuth != nil {
		conf.HttpAuth = sd.consulBasicAuth
	}

	var err error

	sd.ConsulClient, err = consul.NewClient(&conf)
	if err != nil {
		return err
	}

	return nil
}

type ResolvedService struct {
	ID      string
	Address string
	Port    int
}

func (sd *ServiceDiscovery) GetServices(servicename, tag string) ([]ResolvedService, error) {
	options := &consul.QueryOptions{}

	services, _, resolveerr := sd.ConsulClient.Health().Service(servicename, tag, false, options)
	if resolveerr != nil {
		return nil, resolveerr
	}

	var resolvedservices []ResolvedService

	for _, service := range services {
		state := service.Checks.AggregatedStatus()
		if serviceStateIsAcceptable(state) {
			resolvedservices = append(resolvedservices, ResolvedService{
				ID:      service.Service.ID,
				Address: service.Service.Address,
				Port:    service.Service.Port,
			})
		}
	}

	return resolvedservices, nil
}

func serviceStateIsAcceptable(state string) bool {
	return state == consul.HealthPassing || state == consul.HealthWarning
}

func (sd *ServiceDiscovery) GetValue(key string) ([]byte, error) {
	kvpair, _, err := sd.ConsulClient.KV().Get(key, nil)

	if kvpair == nil {
		return []byte{}, err
	}

	return kvpair.Value, err
}

func (sd *ServiceDiscovery) PutValue(key string, value []byte) error {
	_, err := sd.ConsulClient.KV().Put(&consul.KVPair{Key: key, Value: value}, nil)

	return err
}

func (sd *ServiceDiscovery) SubscribeToKey(key string, valuechan chan []byte) {
	actualIndex := uint64(0)

	go func() {
		for {
			kvpair, _, err := sd.ConsulClient.KV().Get(key, &consul.QueryOptions{WaitIndex: actualIndex})
			if err != nil {
				logger.ErrorErrLog(err)
			}

			valuechan <- kvpair.Value

			if actualIndex > kvpair.ModifyIndex {
				actualIndex = 0

				continue
			}

			actualIndex = kvpair.ModifyIndex
		}
	}()
}

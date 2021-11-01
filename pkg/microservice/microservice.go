package microservice

import (
	"fmt"
	"os"

	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/configuration"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/messaging"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/metrics"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/servicediscovery"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/utils"
)

type Microservice struct {
	Name    string
	ID      string
	Version string
	Tags    []string
	Meta    map[string]string

	ServiceURL string

	SDRegistration   *servicediscovery.ServiceRegistration
	ServiceDiscovery *servicediscovery.ServiceDiscovery

	MessageBroker messaging.MessageBroker

	Config *configuration.Configurationobject

	SignalChannel chan os.Signal
}

func (ms *Microservice) StartService() (err error) {
	servicename := ms.getCleanedServiceName()
	err = ms.generateServiceID()

	if err != nil {
		return err
	}

	serviceidentifier := generateServiceIdentifier(servicename, ms)

	ms.logServiceStart()

	ms.initiateReferencedObjects()

	logger.InitGlobalLogger(ms.SignalChannel)

	err = ms.Config.InitGlobalConfiguration(servicename)
	if err != nil {
		return err
	}

	ms.loadMicroserviceSettings()

	ms.MessageBroker.StartMessageBrokerConnection(serviceidentifier)

	ms.prepareServiceRegistration()
	ms.ServiceDiscovery, err = servicediscovery.StartServiceDiscoveryAndRegisterService(ms.SDRegistration)

	if err != nil {
		return err
	}

	return nil
}

func (ms *Microservice) logServiceStart() {
	servicename := ms.getCleanedServiceName()
	serviceid := ms.ID
	servicepid := os.Getpid()
	hostname, err := os.Hostname()
	if err != nil {
		logger.ErrorErrLog(err)
	}

	formatstr := "started service: %s | id: %s | pid: %d | hostname: %s"

	output := fmt.Sprintf(formatstr, servicename, serviceid, servicepid, hostname)

	logger.InfoLog(output)
}

func generateServiceIdentifier(servicename string, ms *Microservice) string {
	serviceidentifier := servicename + "|" + ms.ID

	return serviceidentifier
}

func (ms *Microservice) initiateReferencedObjects() {
	const signalchannelsize = 2

	ms.SDRegistration = &servicediscovery.ServiceRegistration{}
	ms.ServiceDiscovery = &servicediscovery.ServiceDiscovery{}
	ms.Config = &configuration.Configurationobject{}
	ms.SignalChannel = make(chan os.Signal, signalchannelsize)
}

func (ms *Microservice) loadMicroserviceSettings() {
	ms.ServiceURL = ms.Config.GetStringSetting("Service.URL")
	ms.Tags = ms.Config.GetStringSliceSetting("Service.Tags")
	ms.Tags = append(ms.Tags, ms.Version)
	ms.Meta = ms.Config.GetStringMapSettings("ServiceMetaData")

	ms.SDRegistration.ServiceDiscoveryURL = ms.Config.GetStringSetting("ServiceDiscovery.URL")
	ms.SDRegistration.ServiceDiscoveryHealthPingIntervall =
		ms.Config.GetTimeDuration("ServiceDiscovery.HealthPingIntervall")

	ms.MessageBroker.URL = ms.Config.GetStringSetting("MessageBroker.URL")

	ms.loadBasicAuthCredentialsSettings()

	ms.checkAndSetDebugLogLevel()

	ms.checkAndStartMetricsEndpoint()
}

func (ms *Microservice) loadBasicAuthCredentialsSettings() {
	ms.SDRegistration.ServiceDiscoveryUsername = ms.Config.GetStringSetting("ServiceDiscovery.Username")
	ms.SDRegistration.ServiceDiscoveryPassword = ms.Config.GetStringSetting("ServiceDiscovery.Password")

	ms.MessageBroker.Username = ms.Config.GetStringSetting("MessageBroker.Username")
	ms.MessageBroker.Password = ms.Config.GetStringSetting("MessageBroker.Password")
}

func (ms *Microservice) checkAndSetDebugLogLevel() {
	debug := ms.Config.GetBoolSetting("Log.DEBUG")

	if debug {
		logger.SetDefaultLogLevel("debug")
	}
}

func (ms *Microservice) checkAndStartMetricsEndpoint() {
	prometheusaddress := ms.Config.GetStringSetting("Prometheus.Address")

	go func() {
		if prometheusaddress != "" {
			err := metrics.StartMetricsExporter(prometheusaddress)
			if err != nil {
				logger.ErrorErrLog(err)
			}
		}
	}()
}

func (ms *Microservice) prepareServiceRegistration() {
	ms.SDRegistration.MicroserviceID = ms.ID
	ms.SDRegistration.MicroserviceName = ms.Name
	ms.SDRegistration.MicroserviceTags = ms.Tags
	ms.SDRegistration.MicroserviceMetadata = ms.Meta
	ms.SDRegistration.MicroserviceURL = ms.ServiceURL
}

func (ms *Microservice) generateServiceID() error {
	uuid, err := utils.GenerateUUID()

	ms.ID = uuid

	return err
}

func (ms *Microservice) CloseMicroservice() error {
	ms.MessageBroker.CloseConnection()

	err := ms.ServiceDiscovery.DeregisterService()

	return err
}

func (ms *Microservice) getCleanedServiceName() string {
	return utils.TrimAndLowerString(ms.Name)
}

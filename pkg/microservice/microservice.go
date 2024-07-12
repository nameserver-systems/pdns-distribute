package microservice

import (
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/configuration"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/messaging"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/metrics"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/utils"
)

type Microservice struct {
	Name    string
	ID      string
	Version string
	Tags    []string
	Meta    map[string]string

	ServiceURL string

	MessageBroker *messaging.MessageBroker

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

	ms.initializeMessageBroker(serviceidentifier)

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

	ms.Config = &configuration.Configurationobject{}
	ms.SignalChannel = make(chan os.Signal, signalchannelsize)
}

func (ms *Microservice) loadMicroserviceSettings() {
	ms.ServiceURL = ms.Config.GetStringSetting("Service.URL")
	ms.Tags = ms.Config.GetStringSliceSetting("Service.Tags")
	ms.Tags = append(ms.Tags, ms.Version)
	ms.Meta = ms.Config.GetStringMapSettings("ServiceMetaData")

	ms.checkAndSetDebugLogLevel()

	ms.checkAndStartMetricsEndpoint()
}

func (ms *Microservice) initializeMessageBroker(serviceIdentifier string) (err error) {
	brokerURL := ms.Config.GetStringSetting("MessageBroker.URL")
	brokerUser := ms.Config.GetStringSetting("MessageBroker.Username")
	brokerPassword := ms.Config.GetStringSetting("MessageBroker.Password")

	ms.MessageBroker, err = messaging.NewMessageBroker(brokerURL, brokerUser, brokerPassword, serviceIdentifier)

	return
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
				logger.FatalErrLog(err)
			}
		}
	}()
}

func (ms *Microservice) generateServiceID() (err error) {
	ms.ID, err = utils.GenerateUUID()

	return
}

func (ms *Microservice) CloseMicroservice() error {
	return ms.MessageBroker.CloseConnection()
}

func (ms *Microservice) getCleanedServiceName() string {
	return utils.TrimAndLowerString(ms.Name)
}

func (ms *Microservice) GetServicePort() (port string, err error) {
	var parsedURL *url.URL

	if parsedURL, err = url.Parse(ms.ServiceURL); err != nil {
		return
	}

	_, port, err = net.SplitHostPort(parsedURL.Host)

	return
}

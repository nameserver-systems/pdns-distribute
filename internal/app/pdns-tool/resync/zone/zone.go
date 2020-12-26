package zone

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/modelevent"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/configuration"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/messaging"
)

func Execute(zonename string) {
	conf := configuration.Configurationobject{}
	mb := messaging.MessageBroker{}

	zonename = ensureTrailingDot(zonename)

	initializeConnections(conf, mb)

	changeevent := getSettings(&mb, conf)

	payload := prepareResyncRequestPayload(zonename)

	mb.Publish(changeevent, payload)

	logger.InfoLog("Successfully published zone resync for zone: " + zonename)

	mb.CloseConnection()
}

func prepareResyncRequestPayload(zonename string) []byte {
	event := modelevent.ZoneChangeEvent{
		Zone:      zonename,
		ChangedAt: time.Now(),
	}

	payload, marshalerr := json.Marshal(event)
	if marshalerr != nil {
		logger.ErrorErrLog(marshalerr)
	}

	return payload
}

func getSettings(mb *messaging.MessageBroker, conf configuration.Configurationobject) string {
	mb.URL = conf.GetStringSetting("MessageBroker.URL")
	changeevent := conf.GetStringSetting("ZoneEventTopics.Mod")

	return changeevent
}

func initializeConnections(conf configuration.Configurationobject, mb messaging.MessageBroker) {
	const configurationtarget = "pdns-health-checker"

	const toolname = "pdns-tool"

	err := conf.InitGlobalConfiguration(configurationtarget)
	if err != nil {
		logger.FatalErrLog(err)
		os.Exit(1)
	}

	mb.StartMessageBrokerConnection(toolname)
}

func ensureTrailingDot(zonename string) (zoneid string) {
	zoneid = zonename

	if !strings.HasSuffix(zonename, ".") {
		zoneid += "."
	}

	return
}

//nolint:gochecknoglobals
package main

import (
	"fmt"
	"runtime"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/eventlistener"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/servutils"
	msframe "github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	goversion = ""
)

func main() {
	setGoBuildVersion()
	printBuildInfo()

	ms := msframe.Microservice{
		Name:    "pdns-secondary-syncer",
		Version: version,
	}

	startService(&ms)

	go startEventListening(&ms)

	servutils.WaitToShutdownServer(&ms, func() {
		eventlistener.StopWorker()
		eventlistener.StopDNSServer()
		closeService(&ms)
	})
}

func setGoBuildVersion() {
	goversion = runtime.Version()
}

func printBuildInfo() {
	output := fmt.Sprintf("build v%v, commit %v, built with %v, built at %v", version, commit, goversion, date)

	logger.InfoLog(output)
}

func startService(ms *msframe.Microservice) {
	err := ms.StartService()
	if err != nil {
		logger.FatalErrLog(err)
	}
}

func closeService(ms *msframe.Microservice) {
	ms.CloseMicroservice()
}

func startEventListening(ms *msframe.Microservice) {
	err := eventlistener.StartEventListenerAndWorker(ms)
	if err != nil {
		logger.FatalErrLog(err)
	}
}

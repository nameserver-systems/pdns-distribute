//nolint:gochecknoglobals
package main

import (
	"fmt"
	"runtime"

	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthorchestrator"
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
		Name:    "pdns-health-checker",
		Version: version,
	}

	startService(&ms)

	go startHealthServices(&ms)

	servutils.WaitToShutdownServer(&ms, func() {
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
	err := ms.CloseMicroservice()
	if err != nil {
		logger.FatalErrLog(err)
	}
}

func startHealthServices(ms *msframe.Microservice) {
	err := healthorchestrator.StartHealthServices(ms)
	if err != nil {
		logger.FatalErrLog(err)
	}
}

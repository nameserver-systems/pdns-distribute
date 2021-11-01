package servutils

import (
	"os"
	"os/signal"
	"syscall"

	msframe "github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

// WaitToShutdownServer Waits for Interrupt Signal to execute a function which handles closing.
func WaitToShutdownServer(ms *msframe.Microservice, closeFunc func()) {
	const exitcode = 1

	signal.Notify(ms.SignalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	receivedsignal := <-ms.SignalChannel

	closeFunc()

	if receivedsignal == syscall.SIGINT {
		logger.InfoLog("Microservice: " + ms.Name + " stopped fatal.")
		os.Exit(exitcode)
	} else {
		logger.InfoLog("Microservice: " + ms.Name + " stopped successfully.")
	}
}

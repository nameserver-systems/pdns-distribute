//nolint:gochecknoglobals
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"

	"github.com/nats-io/nats.go"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/pkg/servutils"
	msframe "gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
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

	metadata := make(map[string]string, 10)
	metadata["Testinfo"] = "Test"

	xd := msframe.Microservice{
		Name:    "ExampleTestService",
		Version: version,
	}
	err := xd.StartService()
	if err != nil {
		logger.FatalErrLog(err)
	}

	go startHealthServiceEndpoint(xd)

	xd.MessageBroker.SubscribeAsync("test.test", waitformessage)

	servutils.WaitToShutdownServer(&xd, func() {
		closeMicroservice(xd)
		fmt.Println("TEST")
	})
}

func setGoBuildVersion() {
	goversion = runtime.Version()
}

func printBuildInfo() {
	output := fmt.Sprintf("build v%v, commit %v, built with %v, built at %v", version, commit, goversion, date)

	logger.InfoLog(output)
}

func closeMicroservice(xd msframe.Microservice) {
	err := xd.CloseMicroservice()
	if err != nil {
		logger.FatalErrLog(err)
	}
}

func waitformessage(msg *nats.Msg) {
	fmt.Println("Subject: ", msg.Subject)
	fmt.Println("Payload: ", string(msg.Data))
}

func startHealthServiceEndpoint(xd msframe.Microservice) {
	http.HandleFunc("/health", handler)

	sdurl, err2 := url.Parse(xd.ServiceURL)

	if err2 != nil {
		log.Fatal(err2)
	}

	err := http.ListenAndServe(sdurl.Host, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	rsp := healthresp{
		Status: "good",
		Time:   "today",
	}
	jsn, err := json.Marshal(rsp)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = fmt.Fprint(w, string(jsn))
}

type healthresp struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

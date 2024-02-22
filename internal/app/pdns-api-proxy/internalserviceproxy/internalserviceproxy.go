package internalserviceproxy

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-api-proxy/certificate"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-api-proxy/config"
	"github.com/nameserver-systems/pdns-distribute/internal/pkg/httputils"
	msframe "github.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

var (
	ms            *msframe.Microservice        //nolint:gochecknoglobals
	serviceconfig *config.ServiceConfiguration //nolint:gochecknoglobals
)

func StartProxy(microservice *msframe.Microservice) error {
	initConfig(microservice)

	stream, err := ms.MessageBroker.CreatePersistentMessageStore("pdns-distribute-event-store", []string{serviceconfig.AddEventTopic, serviceconfig.ChangeEventTopic, serviceconfig.DeleteEventTopic})
	if err != nil {
		return err
	}

	ms.MessageBroker.SetStream(stream)

	err = startHTTPServer(serviceconfig)
	if err != nil {
		return err
	}

	return nil
}

func initConfig(microservice *msframe.Microservice) {
	ms = microservice
	serviceconfig = config.GetConfiguration(microservice)
}

func startHTTPServer(serviceconfig *config.ServiceConfiguration) error {
	const serverreadtimeout = 10 * time.Second

	const serverwritetimeout = 15 * time.Second

	router := getNewRouterWithRoutes()

	serviceurl := serviceconfig.ServiceURL

	serviceaddress, parseerr := httputils.GetHostAndPortFromURL(serviceurl)
	if parseerr != nil {
		logger.FatalErrLog(parseerr)
	}

	hostname, parseerr2 := httputils.GetHostnameFromURL(serviceurl)
	if parseerr2 != nil {
		logger.FatalErrLog(parseerr2)
	}

	certpath, keypath, certerr := certificate.EnsureCertificatePaths(hostname, serviceconfig)
	if certerr != nil {
		return certerr
	}

	server := &http.Server{
		Addr:              serviceaddress,
		Handler:           router,
		ReadTimeout:       serverreadtimeout,
		WriteTimeout:      serverwritetimeout,
		ReadHeaderTimeout: 5 * time.Second,
	}

	err := server.ListenAndServeTLS(certpath, keypath)
	if err != nil {
		return err
	}

	return nil
}

func getNewRouterWithRoutes() *mux.Router {
	router := getNewRouter()
	router.StrictSlash(true)

	registerRoutes(router)

	return router
}

func getNewRouter() *mux.Router {
	return mux.NewRouter()
}

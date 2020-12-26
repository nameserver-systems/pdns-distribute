//nolint:lll
package internalserviceproxy

import (
	"github.com/gorilla/mux"
)

func registerRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/servers/{server_id}/zones", createZoneHandler(defaultProxy)).Methods("POST")
	router.HandleFunc("/api/v1/servers/{server_id}/zones/{zone_id}/rectify", changeZoneHandler(defaultProxy)).Methods("PUT")
	router.HandleFunc("/api/v1/servers/{server_id}/zones/{zone_id}", changeZoneHandler(defaultProxy)).Methods("PUT")
	router.HandleFunc("/api/v1/servers/{server_id}/zones/{zone_id}", changeZoneHandler(defaultProxy)).Methods("PATCH")
	router.HandleFunc("/api/v1/servers/{server_id}/zones/{zone_id}", deleteZoneHandler(defaultProxy)).Methods("DELETE")

	router.HandleFunc("/api/v1/servers/{server_id}/zones/{zone_id}/metadata", changeZoneHandler(defaultProxy)).Methods("POST")
	router.HandleFunc("/api/v1/servers/{server_id}/zones/{zone_id}/metadata/{metadata_kind}", changeZoneHandler(defaultProxy)).Methods("PUT")
	router.HandleFunc("/api/v1/servers/{server_id}/zones/{zone_id}/metadata/{metadata_kind}", changeZoneHandler(defaultProxy)).Methods("DELETE")

	router.HandleFunc("/api/v1/servers/{server_id}/zones/{zone_id}/cryptokeys", changeZoneHandler(defaultProxy)).Methods("POST")

	router.HandleFunc("/api/v1/servers/{server_id}/zones/{zone_id}/cryptokeys/{cryptokey_id}", changeZoneHandler(defaultProxy)).Methods("PUT")
	router.HandleFunc("/api/v1/servers/{server_id}/zones/{zone_id}/cryptokeys/{cryptokey_id}", changeZoneHandler(defaultProxy)).Methods("DELETE")

	router.PathPrefix("/").HandlerFunc(defaultProxy)
}

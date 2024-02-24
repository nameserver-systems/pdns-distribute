package internalserviceproxy

import (
	"net/http"
)

func registerRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/v1/servers/{server_id}/zones", createZoneHandler(defaultProxy))
	router.HandleFunc("PUT /api/v1/servers/{server_id}/zones/{zone_id}/rectify", changeZoneHandler(defaultProxy))
	router.HandleFunc("PUT /api/v1/servers/{server_id}/zones/{zone_id}", changeZoneHandler(defaultProxy))
	router.HandleFunc("PATCH /api/v1/servers/{server_id}/zones/{zone_id}", changeZoneHandler(defaultProxy))
	router.HandleFunc("DELETE /api/v1/servers/{server_id}/zones/{zone_id}", deleteZoneHandler(defaultProxy))

	router.HandleFunc("POST /api/v1/servers/{server_id}/zones/{zone_id}/metadata", changeZoneHandler(defaultProxy))
	router.HandleFunc("PUT /api/v1/servers/{server_id}/zones/{zone_id}/metadata/{metadata_kind}", changeZoneHandler(defaultProxy))
	router.HandleFunc("DELETE /api/v1/servers/{server_id}/zones/{zone_id}/metadata/{metadata_kind}", changeZoneHandler(defaultProxy))

	router.HandleFunc("POST /api/v1/servers/{server_id}/zones/{zone_id}/cryptokeys", changeZoneHandler(defaultProxy))

	router.HandleFunc("PUT /api/v1/servers/{server_id}/zones/{zone_id}/cryptokeys/{cryptokey_id}", changeZoneHandler(defaultProxy))
	router.HandleFunc("DELETE /api/v1/servers/{server_id}/zones/{zone_id}/cryptokeys/{cryptokey_id}", changeZoneHandler(defaultProxy))

	router.HandleFunc("/", defaultProxy)
}

# pdns-distribute

> Repository is located at https://gitlab.com/nameserver-systems/pdns-distribute

The project [nameserver.systems](https://nameserver.systems) / pdns-distribute contains a set of go based microservices
for syncing multiple powerdns nameservers. The original goal was to build the perfect authoritative dns server
infrastructure. Most registrars and ISPs provide terrible DNS servers for their customers. Well known problems are too few servers,
restricted functionality (e.g. record types) and slow updates. This project provides a solution for those
problems and extended functionality with our service [nameserver.systems](https://nameserver.systems).

These microservices provide a fast, secure and easy way for syncing PowerDNS nameservers.

Microservices:
- pdns-zone-provider (provides secondary-syncer zonefiles as a signed zonestring)
- pdns-api-proxy (generates zone events from api interaction with powerdns primary)
- pdns-health-checker (watch the whole system sync state and triggers specific sync)
- pdns-secondary-syncer (cares about the sync of a local powerdns instance)

## Documentation

[Documentation](https://docs.nameserver.systems) and [Documentation Source](./docs)

## Architecture

- 1x primary (internal data management)
    - clients can using the primary through the powerdns api
    - contains zone information
- Nx secondary (public authoritative nameserver)
    - serves zone data
- Nx consul (service discovery)
    - used for discovering active healthy secondaries
    - healthchecks
        - every microservice pings consul in a configurable interval
        - after three times interval without a ping from a microservice, this service will be unhealthy
        - after ten times interval without a ping the service will be deregistered
        - if the systemload is greater than 10 the service will be in a warn state
- Nx nats (message broker)
    - complete communication between microservices will use handled by broker

## Techstack

* written in Go
* Consul for Service Discovery
* NATS as Message Broker for the biggest amount of communication

## Dependencies

* go (>= 1.16)
* podman
* golangci-lint
* goreleaser
* upx
* shellcheck

made 2019 - 2021 with ❤ by linxside

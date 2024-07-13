# pdns-distribute

The project [nameserver.systems](https://nameserver.systems) / pdns-distribute contains a set of go based microservices
for syncing multiple powerdns nameservers. The original goal was to build the perfect authoritative dns server
infrastructure. Most registrars and ISPs provide terrible DNS servers for their customers. Well known problems are too few servers,
restricted functionality (e.g. record types) and slow updates. This project provides a solution for those
problems and extended functionality with our service [nameserver.systems](https://nameserver.systems).

These microservices provide a fast, secure and easy way for syncing PowerDNS nameservers including the following
features.

## Features
* event based architecture (hidden primary, public secondary) for fast and scalable sync
* self healing of outdated secondaries
* message driven 
* most events parallelized
* easy deployment by providing .deb packages and systemd configs for automatic restart
* security by keeping it simple
* security: and encryption provided by nats
* security: secondaries doesn't have dnssec private keys
* full ipv6 support (ipv4 is optional)

* self-healing
* automation

There are actually four microservices. Three running on the primary for providing zone information and the last one
is running on a secondary for syncing zone data.

!!! danger
    Use at your own risk

Microservices:

- pdns-zone-provider (provides secondary-syncer zonefiles as a signed zonestring)
- pdns-api-proxy (generates zone events from api interaction with powerdns primary)
- pdns-health-checker (watch the whole system sync state and triggers specific sync)
- pdns-secondary-syncer (cares about the sync of a local powerdns instance)

## Architecture

- 1x primary (internal data management)
    - clients can using the primary through the powerdns api
    - contains zone information
- Nx secondary (public authoritative nameserver)
    - serves zone data
- Nx nats (message broker)
    - complete communication between microservices will use handled by broker

### Security

Sensitive data for dnssec signing is kept only on the primary server. The signed zone data - without the secret keys - will be
transferred to every secondary server. All microservices will connect to nats for
secure access
to the infrastructure and encryption of server to server connections. This has the advantage of not having to care about certificates
for each microservice.

## Techstack

* written in Go
* NATS as Message Broker for the biggest amount of communication

## Dependencies

* go (>= 1.17)
* podman
* golangci-lint
* goreleaser
* upx
* shellcheck

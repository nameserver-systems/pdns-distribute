# Test
## Local Test and Development Environment
The local test environment consist of one container per pdns-distribute microservice in one dev image and different infrastructure containers. We depend on
`podman-compose` for building and managing the dev environment. The environment may also work on other operating systems than linux, but there is no guarantee.

*Container Overview:*

??? note "nats"
    The message broker service nats exposes the ports 4222 and 8222 to the host. Nats runs in cluster mode. Every service connects to this nats instance.
??? note "nats2"
The message broker service nats exposes the ports 4223 and 8223 to the host. This instance is part of the two node
cluster which is necessary to activate Jetstream and has no other function.
??? note "pdns-primary"
    This container contains a powerdns server with a sqlite3 database as backend. It is the source
    of zonedata for the synchronization process. This container exposes port 8081 for the api and port 5301 for DNS (like port 53).
??? note "pdns-secondary"
    This container is identical to the pdns-primary except the configured ports. It is the destination for the synchronization
    process. This container exposes port 8082 for the api and port 5300 for DNS (like port 53).
??? note "pdns-zone-provider"
    This container contains every pdns-distribute microservice, but runs only the one from title. The configuration for every microservice is included in this image. It exposes only the metrics endpoint under the port 9500.
??? note "pdns-secondary-syncer"
    This container contains every pdns-distribute microservice, but runs only the one from title. The configuration for every microservice is included in this image. It exposes only the metrics endpoint under the port 9503.
??? note "pdns-health-checker"
    This container contains every pdns-distribute microservice, but runs only the one from title. The configuration for every microservice is included in this image. It exposes only the metrics endpoint under the port 9501.
??? note "pdns-api-proxy"
    This container contains every pdns-distribute microservice, but runs only the one from title. The configuration for every microservice is included in this image. It exposes the metrics endpoint under the port 9502 and the proxy listener reachable using port 30000.

### Build Containers
It is intended to build  when something has changed.

```bash
podman-compose build
```

### Start Containers
```bash
podman-compose up -d --force-recreate -t 1
```

### Stop Container
```bash
podman-compose stop -t 1
```
### Running Binaries directly
```bash
cd cmd/<microservice>/
go run main.go
```

## Linting
### Pipeline
```bash
make lint
```
### Most Linters of golangci-lint
```bash
make golangci-all
```

## Regular Tests
Run all go tests.
```cmd
make test
```

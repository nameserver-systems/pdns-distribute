# Test
## Local Test and Development Environment
The local test environment consist of an OCI container for dependent services. The microservices will run directly 
on the development machine and **not** in a container. This enables a faster development with less overhead. The
makefile contains targets for creating, running and deleting containers. The TOML configurations in `configs/`
are adopted to run in this environment. The environment may also work on other operating systems than linux.

*Container Overview:*

??? note "consul"
    The service discovery tool consul exposes the ports 8300, 8500, 8600. Consul runs as a standalone node.
??? note "nats"
    The messsage broker service nats exposes the ports 4222, 6222, 8222. Nats runs as a standalone node.
??? note "pdns-primary"
    This self-build container contains a powerdns server with a sqlite3 database as backend. It is the source
    of zonedata for the synchronisation process. This container exposes port 8081 for the api and port 5353 as the
    default port for the nameserver (like port 53).
??? note "pdns-secondary"
    This self-build container is identical to the pdns-primary. It is the destinaton for the synchronisation
    process. This container exposes port 18081 for the api and port 53535 as the default port for
    the nameserver (like port 53).


### Build + Create + Start Container
It is intended to run it once and then prefer [Start Container](#start-container).

```bash
make start-create-dev-dependencies
```

### Start Container
If the container exist, start them without recreating them.

```bash
make start-dev-dependencies
```

### Stop Container
```bash
make stop-dev-containers
```
### Running Binaries in test environment
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

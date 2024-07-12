# Installation

## Requirements

* Actual linux based os with systemd support (eg. Debian Buster)

NATS and Consul instances can be run on primaries or secondaries. The only limit is, that only one instance of each
should be run on a system.

* Consul cluster (service discovery)
* NATS cluster (message broker)

### Primary

* PowerDNS >= 4.3 (with database backend)

### Secondaries

* PowerDNS >= 4.3 (with database backend)
* dnsdist >= 1.5 (recommended)

## Install latest version
### Primary
```
wget -O pdns-distribute.deb https://repo.nameserver.systems/latest/pdns-distribute-primary_latest_linux_amd64.deb \
 && dpkg -i pdns-distribute.deb
```
### Secondary
```
wget -O pdns-distribute.deb https://repo.nameserver.systems/latest/pdns-distribute-secondary_latest_linux_amd64.deb \
 && dpkg -i pdns-distribute.deb
```

### How to check if updates available ?
A sha256 checksum file is provided to check for updates and integrity.
```
wget -O pdns-distribute_checksums.txt https://repo.nameserver.systems/latest/pdns-distribute_latest_checksums.txt
```

## Install a specific version
The specific version is part of the url path.

### Primary
```
wget -O pdns-distribute.deb https://repo.nameserver.systems/archive/v0.0.250/pdns-distribute-primary_latest_linux_amd64.deb \
 && dpkg -i pdns-distribute.deb
```
### Secondary
```
wget -O pdns-distribute.deb https://repo.nameserver.systems/archive/v0.0.250/pdns-distribute-secondary_latest_linux_amd64.deb \
 && dpkg -i pdns-distribute.deb
```

### Checksums
A sha256 checksum file is also provided for archived releases.
```
wget -O pdns-distribute_checksums.txt https://repo.nameserver.systems/archive/v0.0.250/pdns-distribute_latest_checksums.txt
```

## Automated installation

!!! attention "To be continued ..."

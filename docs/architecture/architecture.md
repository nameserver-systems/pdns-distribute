# Architecture

The architecture of pdns-distribute provides the following properties:

* event-based
* message-driven
* scalable
* self-healing
* secure by simplicity
* automation

It is intended to be fast, stable, secure and automateable.

## Overview

The architecture consists of four parts (primary, secondary, service discovery / consul, message broker / nats).
It is recommend to have one primary and multiple secondary servers. Multiple Primaries can be used if the database is replicated. In case of multiple
primaries it is important to run the health-checker service on only one instance, due to race conditions by running
multiple health-checkers within the same cluster.Nontheless there must be at least one secondary. The service discovery (consul)
should be run on at least three instances. It is recommended to also run the message broker (nats) on at least three
instances. Those instances can run on separate servers or on primaries and secondaries.

![architecture-overview](../img/pdns-distribute-architecture.png)

## Components
### Nats / message broker 
The message broker provides the communication backend for every microservice. Nats was choosen because it is very lightweight and
powerful.

### Consul / service discovery
The service discovery is used to detect active secondaries. When the discovery process is successful, the secondaries are included in health checks. For the
sake of completeness every microservice registers with consul at startup.

## Security
The consul and nats clusters have to use encryption (TLS) between each instance. Encryption (TLS) and authentication 
are required for all connections to consul and nats. 

Local connections (eg. to powerdns) are not encrypted, without instances of consul and nats.

The api-proxy enforces encryption. If no TLS keys are provided, it will create self-signed keys automatically at startup.
The authentication with the powerdns api uses a token within the http(s) header. The api-proxy will copy the given headers for
requests to the powerdns api.

## Sequence diagrams

??? note "Event Sync"
    ```mermaid
        sequenceDiagram
        participant a as [external-client]
        participant b as internal-proxy
        participant c as [PowerDNS API]
        participant e as external-proxy
        participant d as [Message Broker]
        participant f as secondary-syncer
        a->>+b: Zone-Add
        b-xd: {ZONE.ADD}
        d-xf: {ZONE.ADD}
        f->>+d: Get-Zone-Data
        d->>+e: Get-Zone-Data
        e->>+c: Get-Zone-Data
        c-->>-e: Zone-Data
        e-->>-d: Zone-Data
        d-->>-f: Zone-Data
        b->>+c: Zone-Add
        c-->>-b: Zone-Add-Successful
        b-->>-a: Zone-Add-Successful
        Note over b,c: pdns-primary
        Note over d,e: pdns-primary
        Note over f: pdns-secondary
    ```

??? note "Health Check and trigger Sync"
    ```mermaid
        sequenceDiagram
            participant a as health-checker
            participant b as [Primary PowerDNS API]
            participant c as [Message Broker]
            participant d as [Service Discovery]
            participant e as [Secondary PowerDNS]
            loop Every Quarter of an hour
            a->>+b: Get-Zones-Serial
            b-->>-a: All-Zones-Serial
            a->>+d: Get-Active-Secondaries
            d-->>-a: Active-Secondaries
            loop All Secondaries
            loop All Zones
            a->>+e: Get-Zone-Serial
            e-->>-a: Zone-Serial
            end
            end
            a-xc: {Zone.Change}
            end
            Note over a,b: pdns-primary
            Note over c,d: pdns-primary
            Note over e: pdns-secondary
    ```

??? note "Event based health check"
    ```mermaid
        sequenceDiagram
            participant a as health-checker
            participant b as [Primary PowerDNS API]
            participant c as [Message Broker]
            participant d as [Service Discovery]
            participant e as [Secondary PowerDNS]
            c-xa: {Zone.Add} or {Zone.Change}
            a->>+b: Get-Zone-Serial
            b-->>-a: Zone-Serial
            a->>+d: Get-Active-Secondaries
            d-->>-a: Active-Secondaries
            loop All Secondaries
            a->>+e: Get-Zone-Serial
            e-->>-a: Zone-Serial
            end
            a-xc: {Zone.Change}
            Note over a,b: pdns-primary
            Note over c,d: pdns-primary
            Note over e: pdns-secondary
    ```

??? note "RRSIG Signing Check and Sync"
    ```mermaid
        sequenceDiagram
            participant a as signing-checker
            participant b as [Primary PowerDNS API]
            participant c as [Message Broker]
            loop Every half day
            a->>+b: Get-DNSSEC-Zones-Signature-Validity
            b-->>-a: DNSSEC-Zones-Signature-Validity
            loop All phasing out zones
            a->>+b: Renew Signatures
            a-xc: {Zone.Change}
            end
            end
            Note over a,b: pdns-primary
            Note over c: pdns-primary
    ```

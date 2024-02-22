# pdns-health-checker
The health checker watches the whole system sync state and triggers a specific sync of a secondary if a zone is
outdated. Differences will be recognized by comparing SOA serial. It also triggers refresh of DNSSEC signatures and
sets automatically every quarter of an hour nsec3 for DNSSEC zones with unhashed denial of existence (nsec).
Nsec3 prevents zone walking, by hashing zone names.

## Types of health checks / synchronisations
* event sync check: checks zone propagation after a [configurable](../operation/configuration.md) amount of time after a
zone event
* periodical sync check: check and sync the global zone state in configurable intervals
* signing sync: sync dnssec signed zones after refresh of rrsig (signatures)
* periodical nsec3 activation: activate nsec3 for DNSSEC zones if nsec is used

## Workflows
### Workflow 1 / Event Sync Check
1. Listen on zone events (message broker)
2. Wait for a configurable amount of time after a received event
3. Check zone state on all active and healthy secondaries
4. If zone state is not up-to-date trigger a resync of the outdated secondary

### Workflow 2 / Periodical Sync Check
1. Wait for a configurable amount of time
2. Get active secondaries and zones
3. For every secondary
    1. Get zone state for every active zone from secondary
    2. Get diff between primary and secondary zone state
    3. If difference between primary and secondary then trigger resync of zone

### Workflow 3 / Signing Sync
1. Every week (thursday at 01:00 UTC)
    1. Resync of dnssec signed zones

### Workflow 4 / NSEC3 Activation
1. Every quarter of an hour
    1. Get DNSSEC signed zones
    2. If nsec3 is not activated for zone then activate

### Workflow 5 / Get active secondaries and zones in background
1. Wait for a configurable amount of time
    1. Get active secondaries from nats as active consumers
    2. Get active zones from powerdns primary per HTTP API

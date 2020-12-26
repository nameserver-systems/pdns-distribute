# pdns-zone-provider
The zone provider runs on the primary server. It provides signed zonedata from powerdns for the secondary-syncer.

Can be run multiple in instances on the same primary node.

PowerDNS must be configured to allow api access from localhost and (AXFR).

## Workflows
### Workflow / Respond To Request
1. Waiting on zonedata requests as worker on a message broker topic
    1. Retrieve signed zonedata per AXFR from powerdns
    2. Retrieve zonemetadata per HTTP API from powerdns
    3. Respond zonedata to requesting service

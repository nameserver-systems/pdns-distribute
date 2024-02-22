# pdns-secondary-syncer
The secondary syncer runs on every secondary server. It syncs the local powerdns instance based on incoming zone events
and provides zone state information. The secondary syncer exposes an AXFR interface, which must be used by powerdns for
fetching zone data.

## Workflows
### Workflow / Sync
1. Listen for zone events (message broker)
    1. Handle zone event (add, change, delete)

#### Workflow / Handle Zone Events Add & Change
1. Trigger powerdns via API (create a secondary zone, start AXFR zone transfer)
2. Retrieve AXFR request
3. Retrieve zone data per NATS (message broker)
4. Respond to AXFR request

#### Workflow / Handle Zone Event Delete
1. Trigger powerdns via API (delete secondary zone)

### Workflow / Respond Zone State
1. Listen for secondary state request events (message broker)
    1. Respond the secondaries state via the message broker

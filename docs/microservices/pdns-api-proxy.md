# pdns-api-proxy
The api proxy sits in front of the powerdns primary api. It will generate zone events from api interaction with the
primary api and publish them by using the distributed message broker nats. The microservice will automatically create
self signed certificates for the exposed api listener,
if no certificate will be [configured](../operation/configuration.md).

## Workflows
### Main Workflow
1. Listen for API / HTTPS requests
2. Extract event type and zone from request
3. Publish an event through message broker
4. Forward request to powerdns primary api and respond to the querying client

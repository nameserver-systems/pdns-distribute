# Configuration
All microservices will be configured in [TOML](https://github.com/toml-lang/toml) format. All configurations must
be adopted to your requirements. The packages include default configs.

## Primary
### API Proxy

??? info "/etc/pdns-api-proxy/config.toml"
    ```toml
    [Service]
    URL = "http://localhost:30000"
    Cert = "/home/pdns-api-proxy/cert.pem"
    Key = "/home/pdns-api-proxy/key.pem"
    Tags = ["internal", "primary-node", "proxy"]
        [ServiceMetaData]
    
    # very verbose output, journal may rotate often
    [Log]
    DEBUG = true
    
    [Prometheus]
    Address = "localhost:9502" # turn off by set empty string

    [MessageBroker]
    URL = "nats://localhost:4222"
    Username = "" # optional
    Password = "" # optional
    
    [PowerDNS]
    URL = "http://localhost:8081"
    
    [ZoneEventTopics]
    Add = "zone.add"
    Mod = "zone.modified"
    Del = "zone.delete"
    ```

### Health Checker

??? info "/etc/pdns-health-checker/config.toml"
    ```toml
    [Service]
    URL = "http://localhost:30001"
    Tags = ["internal", "primary-node", "health", "check"]
        [ServiceMetaData]
    
    # very verbose output, journal may rotate often
    [Log]
    DEBUG = true
    
    [Prometheus]
    Address = "localhost:9501" # turn off by set empty string
    
    [MessageBroker]
    URL = "nats://localhost:4222"
    Username = "" # optional
    Password = "" # optional
    
    [PowerDNS]
    URL = "http://localhost:8081"
    APIToken = "0000"
    ServerID = "localhost"
    
    [ZoneEventTopics]
    Add = "zone.add"
    Mod = "zone.modified"
    Del = "zone.delete"
    
    [ZoneStateTopics]
    Prefix = "zonestate.>"
    
    [HealthChecks]
    EventCheckWaitTime = "20s"
    ActiveZoneSecondaryRefreshIntervall = "5m"
    PeriodicalCheckIntervall = "15m"
    NSEC3CheckIntervall = "15m"
    ```

### Zone Provider

??? info "/etc/pdns-zone-provider/config.toml"
    ```toml
    [Service]
    URL = "http://localhost:30001"
    Tags = ["external", "primary-node", "proxy"]
        [ServiceMetaData]
    
    # very verbose output, journal may rotate often
    [Log]
    DEBUG = true
    
    [Prometheus]
    Address = "localhost:9500" # turn off by set empty string
    
    [MessageBroker]
    URL = "nats://localhost:4222"
    Username = "" # optional
    Password = "" # optional
    
    [PowerDNS]
    URL = "http://localhost:8081"
    APIToken = "0000"
    ServerID = "localhost"
    AXFRTimeout = "2s"
    AXFRAddress = "127.0.0.1:53"
    
    [ZoneDataTopics]
    Wildcard = "zonedata.>"
    ```

## Secondary
### Secondary Syncer

??? info "/etc/pdns-secondary-syncer/config.toml"
    ```toml
    [Service]
    URL = "http://localhost:30002"
    Tags = ["internal", "secondary-node", "eu-west"]
        [ServiceMetaData]
    
    # very verbose output, journal may rotate often
    [Log]
    DEBUG = true
    
    [Prometheus]
    Address = "localhost:9503" # turn off by set empty string
    
    [MessageBroker]
    URL = "nats://localhost:4222"
    Username = "" # optional
    Password = "" # optional
    
    [PowerDNS]
    URL = "http://localhost:18081"
    APIToken = "0000"
    ServerID = "localhost"
    EventDelay = "0s"
    APIWorker = 4
    
    [AXFRPrimary]
    Address = "127.0.0.1:20102"
    
    [ZoneEventTopics]
    Add = "zone.add"
    Mod = "zone.modified"
    Del = "zone.delete"
    
    [ZoneDataTopics]
    Prefix = "zonedata."
    
    [ZoneStateTopics]
    Prefix = "zonestate.>"
    ```
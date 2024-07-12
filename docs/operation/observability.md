# Observability
## Logging

All microservices log to syslog. By setting the `DEBUG` flag to `True` in [configurations](./configuration.md) verbose
logging for debugging purposes will be activated. For a production environment `DEBUG` should be `False`.

## Monitoring

All microservices provide a Prometheus compatible endpoint for statics. The Prometheus endpoint can be disabled by setting
an empty string for `Address` in [configurations](./configuration.md). The endpoint is not encrypted and has no
authentication. It is recommended to only bind on local connections like shown in the default config.

We recommend Grafana and Prometheus for displaying and monitoring.

### Grafana Dashboard Examples
![primary-dashboard](../img/pdns-distribute-primary-dashboard.png)
![secondary-dashboar](../img/pdns-distribute-secondary-dashboard.png)


### Prometheus Endpoints
#### Primary - API Proxy

??? info "localhost:9502/metrics"
    ```
    # HELP apiproxy_powerdns_api_requests_total The total count of powerdns api requests
    # TYPE apiproxy_powerdns_api_requests_total counter
    apiproxy_powerdns_api_requests_total 0
    # HELP apiproxy_published_zone_events_total The total count of published zone events
    # TYPE apiproxy_published_zone_events_total counter
    apiproxy_published_zone_events_total{event_type="change"} 0
    apiproxy_published_zone_events_total{event_type="create"} 0
    apiproxy_published_zone_events_total{event_type="delete"} 0
    # HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
    # TYPE go_gc_duration_seconds summary
    go_gc_duration_seconds{quantile="0"} 0
    go_gc_duration_seconds{quantile="0.25"} 0
    go_gc_duration_seconds{quantile="0.5"} 0
    go_gc_duration_seconds{quantile="0.75"} 0
    go_gc_duration_seconds{quantile="1"} 0
    go_gc_duration_seconds_sum 0
    go_gc_duration_seconds_count 0
    # HELP go_goroutines Number of goroutines that currently exist.
    # TYPE go_goroutines gauge
    go_goroutines 16
    # HELP go_info Information about the Go environment.
    # TYPE go_info gauge
    go_info{version="go1.15.2"} 1
    # HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
    # TYPE go_memstats_alloc_bytes gauge
    go_memstats_alloc_bytes 1.43868e+06
    # HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
    # TYPE go_memstats_alloc_bytes_total counter
    go_memstats_alloc_bytes_total 1.43868e+06
    # HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
    # TYPE go_memstats_buck_hash_sys_bytes gauge
    go_memstats_buck_hash_sys_bytes 4608
    # HELP go_memstats_frees_total Total number of frees.
    # TYPE go_memstats_frees_total counter
    go_memstats_frees_total 1107
    # HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
    # TYPE go_memstats_gc_cpu_fraction gauge
    go_memstats_gc_cpu_fraction 0
    # HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
    # TYPE go_memstats_gc_sys_bytes gauge
    go_memstats_gc_sys_bytes 4.110168e+06
    # HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
    # TYPE go_memstats_heap_alloc_bytes gauge
    go_memstats_heap_alloc_bytes 1.43868e+06
    # HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
    # TYPE go_memstats_heap_idle_bytes gauge
    go_memstats_heap_idle_bytes 6.3619072e+07
    # HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
    # TYPE go_memstats_heap_inuse_bytes gauge
    go_memstats_heap_inuse_bytes 2.899968e+06
    # HELP go_memstats_heap_objects Number of allocated objects.
    # TYPE go_memstats_heap_objects gauge
    go_memstats_heap_objects 9193
    # HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
    # TYPE go_memstats_heap_released_bytes gauge
    go_memstats_heap_released_bytes 6.3619072e+07
    # HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
    # TYPE go_memstats_heap_sys_bytes gauge
    go_memstats_heap_sys_bytes 6.651904e+07
    # HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
    # TYPE go_memstats_last_gc_time_seconds gauge
    go_memstats_last_gc_time_seconds 0
    # HELP go_memstats_lookups_total Total number of pointer lookups.
    # TYPE go_memstats_lookups_total counter
    go_memstats_lookups_total 0
    # HELP go_memstats_mallocs_total Total number of mallocs.
    # TYPE go_memstats_mallocs_total counter
    go_memstats_mallocs_total 10300
    # HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
    # TYPE go_memstats_mcache_inuse_bytes gauge
    go_memstats_mcache_inuse_bytes 6944
    # HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
    # TYPE go_memstats_mcache_sys_bytes gauge
    go_memstats_mcache_sys_bytes 16384
    # HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
    # TYPE go_memstats_mspan_inuse_bytes gauge
    go_memstats_mspan_inuse_bytes 54672
    # HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
    # TYPE go_memstats_mspan_sys_bytes gauge
    go_memstats_mspan_sys_bytes 65536
    # HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
    # TYPE go_memstats_next_gc_bytes gauge
    go_memstats_next_gc_bytes 4.473924e+06
    # HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
    # TYPE go_memstats_other_sys_bytes gauge
    go_memstats_other_sys_bytes 1.131184e+06
    # HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
    # TYPE go_memstats_stack_inuse_bytes gauge
    go_memstats_stack_inuse_bytes 589824
    # HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
    # TYPE go_memstats_stack_sys_bytes gauge
    go_memstats_stack_sys_bytes 589824
    # HELP go_memstats_sys_bytes Number of bytes obtained from system.
    # TYPE go_memstats_sys_bytes gauge
    go_memstats_sys_bytes 7.2436744e+07
    # HELP go_threads Number of OS threads created.
    # TYPE go_threads gauge
    go_threads 10
    # HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
    # TYPE process_cpu_seconds_total counter
    process_cpu_seconds_total 0.01
    # HELP process_max_fds Maximum number of open file descriptors.
    # TYPE process_max_fds gauge
    process_max_fds 524288
    # HELP process_open_fds Number of open file descriptors.
    # TYPE process_open_fds gauge
    process_open_fds 12
    # HELP process_resident_memory_bytes Resident memory size in bytes.
    # TYPE process_resident_memory_bytes gauge
    process_resident_memory_bytes 1.4114816e+07
    # HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
    # TYPE process_start_time_seconds gauge
    process_start_time_seconds 1.60296530225e+09
    # HELP process_virtual_memory_bytes Virtual memory size in bytes.
    # TYPE process_virtual_memory_bytes gauge
    process_virtual_memory_bytes 1.33838848e+09
    # HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
    # TYPE process_virtual_memory_max_bytes gauge
    process_virtual_memory_max_bytes -1
    # HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
    # TYPE promhttp_metric_handler_requests_in_flight gauge
    promhttp_metric_handler_requests_in_flight 1
    # HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
    # TYPE promhttp_metric_handler_requests_total counter
    promhttp_metric_handler_requests_total{code="200"} 0
    promhttp_metric_handler_requests_total{code="500"} 0
    promhttp_metric_handler_requests_total{code="503"} 0
    ```

#### Primary - Health Checker

??? info "localhost:9501/metrics"
    ```
    # HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
    # TYPE go_gc_duration_seconds summary
    go_gc_duration_seconds{quantile="0"} 0
    go_gc_duration_seconds{quantile="0.25"} 0
    go_gc_duration_seconds{quantile="0.5"} 0
    go_gc_duration_seconds{quantile="0.75"} 0
    go_gc_duration_seconds{quantile="1"} 0
    go_gc_duration_seconds_sum 0
    go_gc_duration_seconds_count 0
    # HELP go_goroutines Number of goroutines that currently exist.
    # TYPE go_goroutines gauge
    go_goroutines 22
    # HELP go_info Information about the Go environment.
    # TYPE go_info gauge
    go_info{version="go1.15.2"} 1
    # HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
    # TYPE go_memstats_alloc_bytes gauge
    go_memstats_alloc_bytes 1.327232e+06
    # HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
    # TYPE go_memstats_alloc_bytes_total counter
    go_memstats_alloc_bytes_total 1.327232e+06
    # HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
    # TYPE go_memstats_buck_hash_sys_bytes gauge
    go_memstats_buck_hash_sys_bytes 1.445273e+06
    # HELP go_memstats_frees_total Total number of frees.
    # TYPE go_memstats_frees_total counter
    go_memstats_frees_total 540
    # HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
    # TYPE go_memstats_gc_cpu_fraction gauge
    go_memstats_gc_cpu_fraction 0
    # HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
    # TYPE go_memstats_gc_sys_bytes gauge
    go_memstats_gc_sys_bytes 4.091664e+06
    # HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
    # TYPE go_memstats_heap_alloc_bytes gauge
    go_memstats_heap_alloc_bytes 1.327232e+06
    # HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
    # TYPE go_memstats_heap_idle_bytes gauge
    go_memstats_heap_idle_bytes 6.361088e+07
    # HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
    # TYPE go_memstats_heap_inuse_bytes gauge
    go_memstats_heap_inuse_bytes 2.90816e+06
    # HELP go_memstats_heap_objects Number of allocated objects.
    # TYPE go_memstats_heap_objects gauge
    go_memstats_heap_objects 7406
    # HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
    # TYPE go_memstats_heap_released_bytes gauge
    go_memstats_heap_released_bytes 6.3578112e+07
    # HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
    # TYPE go_memstats_heap_sys_bytes gauge
    go_memstats_heap_sys_bytes 6.651904e+07
    # HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
    # TYPE go_memstats_last_gc_time_seconds gauge
    go_memstats_last_gc_time_seconds 0
    # HELP go_memstats_lookups_total Total number of pointer lookups.
    # TYPE go_memstats_lookups_total counter
    go_memstats_lookups_total 0
    # HELP go_memstats_mallocs_total Total number of mallocs.
    # TYPE go_memstats_mallocs_total counter
    go_memstats_mallocs_total 7946
    # HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
    # TYPE go_memstats_mcache_inuse_bytes gauge
    go_memstats_mcache_inuse_bytes 6944
    # HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
    # TYPE go_memstats_mcache_sys_bytes gauge
    go_memstats_mcache_sys_bytes 16384
    # HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
    # TYPE go_memstats_mspan_inuse_bytes gauge
    go_memstats_mspan_inuse_bytes 45968
    # HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
    # TYPE go_memstats_mspan_sys_bytes gauge
    go_memstats_mspan_sys_bytes 49152
    # HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
    # TYPE go_memstats_next_gc_bytes gauge
    go_memstats_next_gc_bytes 4.473924e+06
    # HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
    # TYPE go_memstats_other_sys_bytes gauge
    go_memstats_other_sys_bytes 1.165399e+06
    # HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
    # TYPE go_memstats_stack_inuse_bytes gauge
    go_memstats_stack_inuse_bytes 589824
    # HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
    # TYPE go_memstats_stack_sys_bytes gauge
    go_memstats_stack_sys_bytes 589824
    # HELP go_memstats_sys_bytes Number of bytes obtained from system.
    # TYPE go_memstats_sys_bytes gauge
    go_memstats_sys_bytes 7.3876736e+07
    # HELP go_threads Number of OS threads created.
    # TYPE go_threads gauge
    go_threads 9
    # HELP healthchecker_actual_active_secondary_count The actual count of active secondaries
    # TYPE healthchecker_actual_active_secondary_count gauge
    healthchecker_actual_active_secondary_count 0
    # HELP healthchecker_actual_dnssec_primary_zone_count The actual count of dnssec zones on primary
    # TYPE healthchecker_actual_dnssec_primary_zone_count gauge
    healthchecker_actual_dnssec_primary_zone_count 0
    # HELP healthchecker_actual_primary_zone_count The actual count of zones on primary
    # TYPE healthchecker_actual_primary_zone_count gauge
    healthchecker_actual_primary_zone_count 0
    # HELP healthchecker_ensure_nsec3_cyles_total The total count of ensure nse3 cycles for dnssec zones
    # TYPE healthchecker_ensure_nsec3_cyles_total counter
    healthchecker_ensure_nsec3_cyles_total 0
    # HELP healthchecker_event_health_check_published_zone_events_total The total count of published zone events
    # TYPE healthchecker_event_health_check_published_zone_events_total counter
    healthchecker_event_health_check_published_zone_events_total{event_type="change"} 0
    healthchecker_event_health_check_published_zone_events_total{event_type="create"} 0
    healthchecker_event_health_check_published_zone_events_total{event_type="delete"} 0
    # HELP healthchecker_event_health_check_received_zone_events_total The total count of received zone events
    # TYPE healthchecker_event_health_check_received_zone_events_total counter
    healthchecker_event_health_check_received_zone_events_total{event_type="change"} 0
    healthchecker_event_health_check_received_zone_events_total{event_type="create"} 0
    healthchecker_event_health_check_received_zone_events_total{event_type="delete"} 0
    # HELP healthchecker_interval_health_cycles_total The total count of interval health check cycles
    # TYPE healthchecker_interval_health_cycles_total counter
    healthchecker_interval_health_cycles_total 0
    # HELP healthchecker_interval_health_published_zone_events_total The total count of published zone events
    # TYPE healthchecker_interval_health_published_zone_events_total counter
    healthchecker_interval_health_published_zone_events_total{event_type="change"} 0
    healthchecker_interval_health_published_zone_events_total{event_type="create"} 0
    healthchecker_interval_health_published_zone_events_total{event_type="delete"} 0
    # HELP healthchecker_interval_signing_sync_published_zone_events_total The total count of published zone events
    # TYPE healthchecker_interval_signing_sync_published_zone_events_total counter
    healthchecker_interval_signing_sync_published_zone_events_total{event_type="change"} 0
    # HELP healthchecker_powerdns_api_call_total The total count of powerdns api calls
    # TYPE healthchecker_powerdns_api_call_total counter
    healthchecker_powerdns_api_call_total{state="failed"} 0
    healthchecker_powerdns_api_call_total{state="successful"} 2
    # HELP healthchecker_powerdns_cli_call_total The total count of powerdns cli calls
    # TYPE healthchecker_powerdns_cli_call_total counter
    healthchecker_powerdns_cli_call_total{state="failed"} 0
    healthchecker_powerdns_cli_call_total{state="successful"} 0
    # HELP healthchecker_refresh_cycles_total The total count of refresh cycles of primary zones and active secondaries
    # TYPE healthchecker_refresh_cycles_total counter
    healthchecker_refresh_cycles_total 0
    # HELP healthchecker_service_discovery_call_total The total count of service discovery calls
    # TYPE healthchecker_service_discovery_call_total counter
    healthchecker_service_discovery_call_total{state="failed"} 0
    healthchecker_service_discovery_call_total{state="successful"} 1
    # HELP healthchecker_singing_sync_cycles_total The total count of signing sync cycles of dnssec zones
    # TYPE healthchecker_singing_sync_cycles_total counter
    healthchecker_singing_sync_cycles_total 0
    # HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
    # TYPE process_cpu_seconds_total counter
    process_cpu_seconds_total 0.01
    # HELP process_max_fds Maximum number of open file descriptors.
    # TYPE process_max_fds gauge
    process_max_fds 524288
    # HELP process_open_fds Number of open file descriptors.
    # TYPE process_open_fds gauge
    process_open_fds 11
    # HELP process_resident_memory_bytes Resident memory size in bytes.
    # TYPE process_resident_memory_bytes gauge
    process_resident_memory_bytes 1.4700544e+07
    # HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
    # TYPE process_start_time_seconds gauge
    process_start_time_seconds 1.60296544639e+09
    # HELP process_virtual_memory_bytes Virtual memory size in bytes.
    # TYPE process_virtual_memory_bytes gauge
    process_virtual_memory_bytes 1.265332224e+09
    # HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
    # TYPE process_virtual_memory_max_bytes gauge
    process_virtual_memory_max_bytes -1
    # HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
    # TYPE promhttp_metric_handler_requests_in_flight gauge
    promhttp_metric_handler_requests_in_flight 1
    # HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
    # TYPE promhttp_metric_handler_requests_total counter
    promhttp_metric_handler_requests_total{code="200"} 0
    promhttp_metric_handler_requests_total{code="500"} 0
    promhttp_metric_handler_requests_total{code="503"} 0
    ```

#### Primary - Zone Provider

??? info "localhost:9500/metrics"
    ```
    # HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
    # TYPE go_gc_duration_seconds summary
    go_gc_duration_seconds{quantile="0"} 0
    go_gc_duration_seconds{quantile="0.25"} 0
    go_gc_duration_seconds{quantile="0.5"} 0
    go_gc_duration_seconds{quantile="0.75"} 0
    go_gc_duration_seconds{quantile="1"} 0
    go_gc_duration_seconds_sum 0
    go_gc_duration_seconds_count 0
    # HELP go_goroutines Number of goroutines that currently exist.
    # TYPE go_goroutines gauge
    go_goroutines 16
    # HELP go_info Information about the Go environment.
    # TYPE go_info gauge
    go_info{version="go1.15.2"} 1
    # HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
    # TYPE go_memstats_alloc_bytes gauge
    go_memstats_alloc_bytes 1.212408e+06
    # HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
    # TYPE go_memstats_alloc_bytes_total counter
    go_memstats_alloc_bytes_total 1.212408e+06
    # HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
    # TYPE go_memstats_buck_hash_sys_bytes gauge
    go_memstats_buck_hash_sys_bytes 4629
    # HELP go_memstats_frees_total Total number of frees.
    # TYPE go_memstats_frees_total counter
    go_memstats_frees_total 510
    # HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
    # TYPE go_memstats_gc_cpu_fraction gauge
    go_memstats_gc_cpu_fraction 0
    # HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
    # TYPE go_memstats_gc_sys_bytes gauge
    go_memstats_gc_sys_bytes 4.091664e+06
    # HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
    # TYPE go_memstats_heap_alloc_bytes gauge
    go_memstats_heap_alloc_bytes 1.212408e+06
    # HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
    # TYPE go_memstats_heap_idle_bytes gauge
    go_memstats_heap_idle_bytes 6.3823872e+07
    # HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
    # TYPE go_memstats_heap_inuse_bytes gauge
    go_memstats_heap_inuse_bytes 2.727936e+06
    # HELP go_memstats_heap_objects Number of allocated objects.
    # TYPE go_memstats_heap_objects gauge
    go_memstats_heap_objects 6637
    # HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
    # TYPE go_memstats_heap_released_bytes gauge
    go_memstats_heap_released_bytes 6.3791104e+07
    # HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
    # TYPE go_memstats_heap_sys_bytes gauge
    go_memstats_heap_sys_bytes 6.6551808e+07
    # HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
    # TYPE go_memstats_last_gc_time_seconds gauge
    go_memstats_last_gc_time_seconds 0
    # HELP go_memstats_lookups_total Total number of pointer lookups.
    # TYPE go_memstats_lookups_total counter
    go_memstats_lookups_total 0
    # HELP go_memstats_mallocs_total Total number of mallocs.
    # TYPE go_memstats_mallocs_total counter
    go_memstats_mallocs_total 7147
    # HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
    # TYPE go_memstats_mcache_inuse_bytes gauge
    go_memstats_mcache_inuse_bytes 6944
    # HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
    # TYPE go_memstats_mcache_sys_bytes gauge
    go_memstats_mcache_sys_bytes 16384
    # HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
    # TYPE go_memstats_mspan_inuse_bytes gauge
    go_memstats_mspan_inuse_bytes 54672
    # HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
    # TYPE go_memstats_mspan_sys_bytes gauge
    go_memstats_mspan_sys_bytes 65536
    # HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
    # TYPE go_memstats_next_gc_bytes gauge
    go_memstats_next_gc_bytes 4.473924e+06
    # HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
    # TYPE go_memstats_other_sys_bytes gauge
    go_memstats_other_sys_bytes 887523
    # HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
    # TYPE go_memstats_stack_inuse_bytes gauge
    go_memstats_stack_inuse_bytes 557056
    # HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
    # TYPE go_memstats_stack_sys_bytes gauge
    go_memstats_stack_sys_bytes 557056
    # HELP go_memstats_sys_bytes Number of bytes obtained from system.
    # TYPE go_memstats_sys_bytes gauge
    go_memstats_sys_bytes 7.21746e+07
    # HELP go_threads Number of OS threads created.
    # TYPE go_threads gauge
    go_threads 10
    # HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
    # TYPE process_cpu_seconds_total counter
    process_cpu_seconds_total 0.01
    # HELP process_max_fds Maximum number of open file descriptors.
    # TYPE process_max_fds gauge
    process_max_fds 524288
    # HELP process_open_fds Number of open file descriptors.
    # TYPE process_open_fds gauge
    process_open_fds 11
    # HELP process_resident_memory_bytes Resident memory size in bytes.
    # TYPE process_resident_memory_bytes gauge
    process_resident_memory_bytes 1.3283328e+07
    # HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
    # TYPE process_start_time_seconds gauge
    process_start_time_seconds 1.60296554581e+09
    # HELP process_virtual_memory_bytes Virtual memory size in bytes.
    # TYPE process_virtual_memory_bytes gauge
    process_virtual_memory_bytes 1.339047936e+09
    # HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
    # TYPE process_virtual_memory_max_bytes gauge
    process_virtual_memory_max_bytes -1
    # HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
    # TYPE promhttp_metric_handler_requests_in_flight gauge
    promhttp_metric_handler_requests_in_flight 1
    # HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
    # TYPE promhttp_metric_handler_requests_total counter
    promhttp_metric_handler_requests_total{code="200"} 0
    promhttp_metric_handler_requests_total{code="500"} 0
    promhttp_metric_handler_requests_total{code="503"} 0
    # HELP zoneprovider_powerdns_api_call_total The total count of powerdns api calls
    # TYPE zoneprovider_powerdns_api_call_total counter
    zoneprovider_powerdns_api_call_total{state="failed"} 0
    zoneprovider_powerdns_api_call_total{state="successful"} 0
    # HELP zoneprovider_powerdns_dns_call_total The total count of powerdns dns calls
    # TYPE zoneprovider_powerdns_dns_call_total counter
    zoneprovider_powerdns_dns_call_total{protocol="tcp",state="failed"} 0
    zoneprovider_powerdns_dns_call_total{protocol="tcp",state="successful"} 0
    # HELP zoneprovider_powerdns_dns_query_total The total count of powerdns dns queries
    # TYPE zoneprovider_powerdns_dns_query_total counter
    zoneprovider_powerdns_dns_query_total{qtype="axfr",state="failed"} 0
    zoneprovider_powerdns_dns_query_total{qtype="axfr",state="successful"} 0
    zoneprovider_powerdns_dns_query_total{qtype="soa",state="failed"} 0
    zoneprovider_powerdns_dns_query_total{qtype="soa",state="successful"} 0
    # HELP zoneprovider_request_total The total count of requests
    # TYPE zoneprovider_request_total counter
    zoneprovider_request_total 0
    # HELP zoneprovider_response_total The total count of responses
    # TYPE zoneprovider_response_total counter
    zoneprovider_response_total 0
    ```
    
#### Secondary - Secondary Syncer

??? info "localhost:9503/metrics"
    ```
    # HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
    # TYPE go_gc_duration_seconds summary
    go_gc_duration_seconds{quantile="0"} 0
    go_gc_duration_seconds{quantile="0.25"} 0
    go_gc_duration_seconds{quantile="0.5"} 0
    go_gc_duration_seconds{quantile="0.75"} 0
    go_gc_duration_seconds{quantile="1"} 0
    go_gc_duration_seconds_sum 0
    go_gc_duration_seconds_count 0
    # HELP go_goroutines Number of goroutines that currently exist.
    # TYPE go_goroutines gauge
    go_goroutines 29
    # HELP go_info Information about the Go environment.
    # TYPE go_info gauge
    go_info{version="go1.15.2"} 1
    # HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
    # TYPE go_memstats_alloc_bytes gauge
    go_memstats_alloc_bytes 1.25928e+06
    # HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
    # TYPE go_memstats_alloc_bytes_total counter
    go_memstats_alloc_bytes_total 1.25928e+06
    # HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
    # TYPE go_memstats_buck_hash_sys_bytes gauge
    go_memstats_buck_hash_sys_bytes 1.444886e+06
    # HELP go_memstats_frees_total Total number of frees.
    # TYPE go_memstats_frees_total counter
    go_memstats_frees_total 529
    # HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
    # TYPE go_memstats_gc_cpu_fraction gauge
    go_memstats_gc_cpu_fraction 0
    # HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
    # TYPE go_memstats_gc_sys_bytes gauge
    go_memstats_gc_sys_bytes 4.097832e+06
    # HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
    # TYPE go_memstats_heap_alloc_bytes gauge
    go_memstats_heap_alloc_bytes 1.25928e+06
    # HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
    # TYPE go_memstats_heap_idle_bytes gauge
    go_memstats_heap_idle_bytes 6.3782912e+07
    # HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
    # TYPE go_memstats_heap_inuse_bytes gauge
    go_memstats_heap_inuse_bytes 2.768896e+06
    # HELP go_memstats_heap_objects Number of allocated objects.
    # TYPE go_memstats_heap_objects gauge
    go_memstats_heap_objects 7001
    # HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
    # TYPE go_memstats_heap_released_bytes gauge
    go_memstats_heap_released_bytes 6.3750144e+07
    # HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
    # TYPE go_memstats_heap_sys_bytes gauge
    go_memstats_heap_sys_bytes 6.6551808e+07
    # HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
    # TYPE go_memstats_last_gc_time_seconds gauge
    go_memstats_last_gc_time_seconds 0
    # HELP go_memstats_lookups_total Total number of pointer lookups.
    # TYPE go_memstats_lookups_total counter
    go_memstats_lookups_total 0
    # HELP go_memstats_mallocs_total Total number of mallocs.
    # TYPE go_memstats_mallocs_total counter
    go_memstats_mallocs_total 7530
    # HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
    # TYPE go_memstats_mcache_inuse_bytes gauge
    go_memstats_mcache_inuse_bytes 6944
    # HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
    # TYPE go_memstats_mcache_sys_bytes gauge
    go_memstats_mcache_sys_bytes 16384
    # HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
    # TYPE go_memstats_mspan_inuse_bytes gauge
    go_memstats_mspan_inuse_bytes 54808
    # HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
    # TYPE go_memstats_mspan_sys_bytes gauge
    go_memstats_mspan_sys_bytes 65536
    # HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
    # TYPE go_memstats_next_gc_bytes gauge
    go_memstats_next_gc_bytes 4.473924e+06
    # HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
    # TYPE go_memstats_other_sys_bytes gauge
    go_memstats_other_sys_bytes 1.143234e+06
    # HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
    # TYPE go_memstats_stack_inuse_bytes gauge
    go_memstats_stack_inuse_bytes 557056
    # HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
    # TYPE go_memstats_stack_sys_bytes gauge
    go_memstats_stack_sys_bytes 557056
    # HELP go_memstats_sys_bytes Number of bytes obtained from system.
    # TYPE go_memstats_sys_bytes gauge
    go_memstats_sys_bytes 7.3876736e+07
    # HELP go_threads Number of OS threads created.
    # TYPE go_threads gauge
    go_threads 10
    # HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
    # TYPE process_cpu_seconds_total counter
    process_cpu_seconds_total 0
    # HELP process_max_fds Maximum number of open file descriptors.
    # TYPE process_max_fds gauge
    process_max_fds 524288
    # HELP process_open_fds Number of open file descriptors.
    # TYPE process_open_fds gauge
    process_open_fds 13
    # HELP process_resident_memory_bytes Resident memory size in bytes.
    # TYPE process_resident_memory_bytes gauge
    process_resident_memory_bytes 1.4393344e+07
    # HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
    # TYPE process_start_time_seconds gauge
    process_start_time_seconds 1.60296564834e+09
    # HELP process_virtual_memory_bytes Virtual memory size in bytes.
    # TYPE process_virtual_memory_bytes gauge
    process_virtual_memory_bytes 1.274044416e+09
    # HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
    # TYPE process_virtual_memory_max_bytes gauge
    process_virtual_memory_max_bytes -1
    # HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
    # TYPE promhttp_metric_handler_requests_in_flight gauge
    promhttp_metric_handler_requests_in_flight 1
    # HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
    # TYPE promhttp_metric_handler_requests_total counter
    promhttp_metric_handler_requests_total{code="200"} 0
    promhttp_metric_handler_requests_total{code="500"} 0
    promhttp_metric_handler_requests_total{code="503"} 0
    # HELP secondarysyncer_local_dns_proxy_queries_total The total count of local dns queries
    # TYPE secondarysyncer_local_dns_proxy_queries_total counter
    secondarysyncer_local_dns_proxy_queries_total{qtype="*",state="prohibited"} 0
    secondarysyncer_local_dns_proxy_queries_total{qtype="axfr",state="permitted"} 0
    secondarysyncer_local_dns_proxy_queries_total{qtype="soa",state="permitted"} 0
    # HELP secondarysyncer_powerdns_api_call_total The total count of powerdns api calls
    # TYPE secondarysyncer_powerdns_api_call_total counter
    secondarysyncer_powerdns_api_call_total{state="failed",type="create"} 0
    secondarysyncer_powerdns_api_call_total{state="failed",type="delete"} 0
    secondarysyncer_powerdns_api_call_total{state="failed",type="get_zones"} 0
    secondarysyncer_powerdns_api_call_total{state="successful",type="clear_cache"} 0
    secondarysyncer_powerdns_api_call_total{state="successful",type="create"} 0
    secondarysyncer_powerdns_api_call_total{state="successful",type="delete"} 0
    secondarysyncer_powerdns_api_call_total{state="successful",type="get_zones"} 0
    # HELP secondarysyncer_powerdns_cli_call_total The total count of powerdns cli calls
    # TYPE secondarysyncer_powerdns_cli_call_total counter
    secondarysyncer_powerdns_cli_call_total{state="failed"} 0
    secondarysyncer_powerdns_cli_call_total{state="successful"} 0
    # HELP secondarysyncer_received_zone_events_total The total count of received zone events
    # TYPE secondarysyncer_received_zone_events_total counter
    secondarysyncer_received_zone_events_total{event_type="change"} 0
    secondarysyncer_received_zone_events_total{event_type="create"} 0
    secondarysyncer_received_zone_events_total{event_type="delete"} 0
    # HELP secondarysyncer_worker_lock_count The actual count of worker locks
    # TYPE secondarysyncer_worker_lock_count gauge
    secondarysyncer_worker_lock_count 0
    # HELP secondarysyncer_zone_data_request_total The total count of zone data requests
    # TYPE secondarysyncer_zone_data_request_total counter
    secondarysyncer_zone_data_request_total{state="failed"} 0
    secondarysyncer_zone_data_request_total{state="successful"} 0
    # HELP secondarysyncer_zone_state_request_total The total count of zone state requests
    # TYPE secondarysyncer_zone_state_request_total counter
    secondarysyncer_zone_state_request_total 0
    # HELP secondarysyncer_zone_state_response_total The total count of zone state responses
    # TYPE secondarysyncer_zone_state_response_total counter
    secondarysyncer_zone_state_response_total{state="failed"} 0
    secondarysyncer_zone_state_response_total{state="successful"} 0
    ```

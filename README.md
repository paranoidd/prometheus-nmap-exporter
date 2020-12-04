# prometheus-nmap-exporter

This code is a kind of PoC regarding dead-simple external firewall monitoring.
For functional monitoring of TCP listeners, look at [prometheus/blackbox_exporter](https://github.com/prometheus/blackbox_exporter)

This also served as Golang learning experience.

### So what happens here?
This code will:
1. read targets from config.json file (see `config.json.example`)
2. Run simple nmap scan against these targets
3. Compare the scan results with input configuration and display appropriate metrics

### Metrics

|Metric name|Description|Values|
|-------------|-------------|-----|
|service_monitor_port_expected|Describes if discovered port was on the expected list|0 for unexpected host/port combination, 1 for expected|
|service_monitor_port_missing|Describes if a port that was not the expected list was not found|0 for present port, 1 for missing port|

### Metrics Labels
Two labels are attached to each metric:
- host - with the name of scanned host
- port - which port is that particular metric value about

### TODO
- displayed metrics need descriptions, type definitions etc (could use own struct)
- make configurable: WithSYNScan() / WithConnectScan() (SYN scan tends to be very slow in docker)
- UDP scans?
- comments + docs
- tests
- automated build / pipeline using public CI systems
- docker hub(?)

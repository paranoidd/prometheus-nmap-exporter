package main

// All const used here should be configurable via ENV variable; omitting the _DEFAULT suffix from it's name

const NMAP_EXPORTER_APP_CONFIG_FILE_DEFAULT = "config.json"
const NMAP_EXPORTER_BIND_PORT_DEFAULT = "9777"
const NMAP_EXPORTER_BIND_SERVER_DEFAULT = "127.0.0.1"
const NMAP_EXPORTER_LOG_REQUESTS_DEFAULT = "true"
const NMAP_EXPORTER_METRICS_FILE_DEFAULT = "/tmp/nmap-exporter-metrics.txt"
const NMAP_EXPORTER_SCAN_INTERVAL_DEFAULT = "30"

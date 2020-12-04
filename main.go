package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var metricsFile string

func main() {
	// Not sure if this is not stupid, but I hope somebody will point it out if it is
	metricsFile = os.Getenv("NMAP_EXPORTER_METRICS_FILE")

	if metricsFile == "" {
		log.Printf("Metrics output file (env NMAP_EXPORTER_METRICS_FILE) not configured, using %s", NMAP_EXPORTER_METRICS_FILE_DEFAULT)
		metricsFile = NMAP_EXPORTER_METRICS_FILE_DEFAULT
	}

	// Remove any potential stale metrics files
	err := os.RemoveAll(metricsFile)
	check(err)

	// We will use two goroutines here
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		intervalStr := os.Getenv("NMAP_EXPORTER_SCAN_INTERVAL")

		if intervalStr == "" {
			log.Printf("Nmap scan interval (env NMAP_EXPORTER_SCAN_INTERVAL) not configured, using %s", NMAP_EXPORTER_SCAN_INTERVAL_DEFAULT)
			intervalStr = NMAP_EXPORTER_SCAN_INTERVAL_DEFAULT
		}

		interval, err := strconv.ParseInt(intervalStr, 10, 64)
		check(err)

		appConfigFile := os.Getenv("NMAP_EXPORTER_APP_CONFIG_FILE")

		if appConfigFile == "" {
			log.Printf("Nmap hosts config file (env NMAP_EXPORTER_APP_CONFIG_FILE) not configured, using %s", NMAP_EXPORTER_APP_CONFIG_FILE_DEFAULT)
			appConfigFile = NMAP_EXPORTER_APP_CONFIG_FILE_DEFAULT
		}

		appConfig, err := ioutil.ReadFile(appConfigFile)
		check(err)

		var expectedList []Host
		json.Unmarshal([]byte(appConfig), &expectedList)

		for {
			NmapScan(expectedList)
			time.Sleep(time.Duration(interval) * time.Second)
		}
		wg.Done()
	}()

	go func() {
		ServerListener()
		wg.Done()
	}()

	// Wait for goroutines to exit
	wg.Wait()
}

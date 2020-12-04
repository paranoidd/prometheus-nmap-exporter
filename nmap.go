package main

import (
	"context"
	"fmt"
	"github.com/Ullaakut/nmap/v2"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type Host struct {
	Name  string
	Ports []string
}

type Result struct {
	Hosts []HostResult
}

type HostResult struct {
	Name  string
	Ports []PortResult
}

type PortResult struct {
	Port     string
	Expected bool
	Missing  bool
}

func NmapScan(expectedList []Host) {

	var discoveredList []Host
	var targets []string

	log.Println("Running nmap scan")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for i := range expectedList {
		targets = append(targets, expectedList[i].Name)
	}

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(targets...),
		nmap.WithContext(ctx),
		nmap.WithConnectScan(),
	)
	check(err)

	nmapResult, warnings, err := scanner.Run()
	check(err)

	log.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(nmapResult.Hosts), nmapResult.Stats.Finished.Elapsed)

	if warnings != nil {
		log.Printf("Nmap warnings: \n %v", warnings)
	}

	discoveredList = GetDiscoveredList(nmapResult)

	WriteMetricsFile(CompareResults(expectedList, discoveredList))
}

//
func GetDiscoveredList(nmapResult *nmap.Run) []Host {
	var discoveredList []Host

	// Use the results to print an example output
	for _, host := range nmapResult.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		h := Host{
			Name: strings.Replace(fmt.Sprintf("%q", host.Hostnames[0]), "\"", "", -1),
		}

		for _, port := range host.Ports {
			h.Ports = append(h.Ports, fmt.Sprintf("%d/%s", port.ID, port.Protocol))
		}

		discoveredList = append(discoveredList, h)
	}

	return discoveredList
}

// Compare nmap scan results with expected configuration
func CompareResults(expected []Host, discovered []Host) Result {
	var r Result

	// Check if all hosts scanned as expected
	for e := range expected {
		// Empty Name means host was not in nmap scan results
		if getHostByName(expected[e].Name, discovered).Name == "" {
			h := HostResult{
				Name: expected[e].Name,
			}

			// All expected ports will be marked as missing due to host not being present in scan result
			for p := range expected[e].Ports {
				h.Ports = append(h.Ports, PortResult{
					Port:     expected[e].Ports[p],
					Expected: true,
					Missing:  true,
				})
			}

			r.Hosts = append(r.Hosts, h)
		}
	}

	// Check the discoveries
	for d := range discovered {
		exp := getHostByName(discovered[d].Name, expected)

		if exp.Name == "" {
			log.Fatalf("Discovered host %s not on the expected list!", discovered[d].Name)
		}

		h := HostResult{
			Name: exp.Name,
		}

		// First loop over expected ports and check whether they are in the scan
		for p := range exp.Ports {
			h.Ports = append(h.Ports, PortResult{
				Port:     exp.Ports[p],
				Expected: true,
				Missing:  !stringInSlice(exp.Ports[p], discovered[d].Ports),
			})
		}

		// Next loop over discovered ports
		for p := range discovered[d].Ports {
			// Skip ports that are on expected list already (no duplications)
			if !stringInSlice(discovered[d].Ports[p], exp.Ports) {
				h.Ports = append(h.Ports, PortResult{
					Port:     discovered[d].Ports[p],
					Expected: false,
					Missing:  false,
				})
			}
		}

		r.Hosts = append(r.Hosts, h)
	}

	return r
}

//
func GetMetricsFromResult(r Result) []string {
	var metrics []string

	for i := range r.Hosts {
		for p := range r.Hosts[i].Ports {
			var portMissing int8
			var portExpected int8

			if r.Hosts[i].Ports[p].Missing {
				portMissing = 1
			}

			if r.Hosts[i].Ports[p].Expected {
				portExpected = 1
			}

			metrics = append(metrics, fmt.Sprintf("service_monitor_port_expected{host=\"%s\",port=\"%s\"} %d\n", r.Hosts[i].Name, r.Hosts[i].Ports[p].Port, portExpected))
			metrics = append(metrics, fmt.Sprintf("service_monitor_port_missing{host=\"%s\",port=\"%s\"} %d\n", r.Hosts[i].Name, r.Hosts[i].Ports[p].Port, portMissing))
		}
	}

	sort.Sort(sort.StringSlice(metrics))

	return metrics
}

// This will ultimately be a route in server
func WriteMetricsFile(r Result) {
	f, err := os.Create(metricsFile)
	check(err)
	defer f.Close()

	metrics := GetMetricsFromResult(r)

	for m := range metrics {
		_, err := f.WriteString(metrics[m])
		check(err)
	}

	f.Sync()

	log.Printf("Wrote new metrics to %s", metricsFile)
}

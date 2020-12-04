package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Log requests to API server
func ServerRequestLogger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t\t%s\t\t%s\t\t%s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

func ServerListener() {
	// Server config
	bindServer := os.Getenv("NMAP_EXPORTER_BIND_SERVER")
	bindPort := os.Getenv("NMAP_EXPORTER_BIND_PORT")

	// Use some defaults if nothing is configured
	if bindServer == "" {
		log.Printf("Bind server (env NMAP_EXPORTER_BIND_SERVER) not configured, using %s", NMAP_EXPORTER_BIND_SERVER_DEFAULT)
		bindServer = NMAP_EXPORTER_BIND_SERVER_DEFAULT
	}

	if bindPort == "" {
		log.Printf("Bind port (env NMAP_EXPORTER_BIND_PORT) not configured, using %s", NMAP_EXPORTER_BIND_PORT_DEFAULT)
		bindPort = NMAP_EXPORTER_BIND_PORT_DEFAULT
	}

	// Routing
	r := mux.NewRouter()

	// Check if logging requests should be enabled
	logReqStr := os.Getenv("NMAP_EXPORTER_LOG_REQUESTS")

	if logReqStr == "" {
		log.Printf("Logging HTTP requests (env NMAP_EXPORTER_LOG_REQUESTS) not configured, using %s", NMAP_EXPORTER_LOG_REQUESTS_DEFAULT)
		logReqStr = NMAP_EXPORTER_LOG_REQUESTS_DEFAULT
	}

	logReqBool, err := strconv.ParseBool(logReqStr)
	check(err)

	if logReqBool {
		// Log requests
		r.Use(ServerRequestLogger)
	}

	// Route to display metrics
	r.HandleFunc("/metrics", RouteMetrics).Methods(http.MethodGet)

	log.Printf("Staring server on %s:%s", bindServer, bindPort)
	log.Fatal(http.ListenAndServe(bindServer+":"+bindPort, r))
}

func RouteMetrics(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, metricsFile)
}

package main

import (
	"flag"
	"fmt"
	exporter "github.com/ovnworks/gobgp_exporter/pkg/gobgp_exporter"
	"github.com/prometheus/common/log"
	"net/http"
	"os"
)

func main() {
	var listenAddress string
	var metricsPath string
	var serverAddress string
	var pollTimeout int
	var pollInterval int
	var isShowMetrics bool
	var isShowVersion bool
	var logLevel string
	var authToken string

	flag.StringVar(&listenAddress, "web.listen-address", ":9474", "Address to listen on for web interface and telemetry.")
	flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	flag.StringVar(&serverAddress, "gobgp.address", "127.0.0.1:50051", "gRPC API address of GoBGP server.")
	flag.IntVar(&pollTimeout, "gobgp.timeout", 2, "Timeout on gRPC requests to a GoBGP server.")
	flag.IntVar(&pollInterval, "gobgp.poll-interval", 15, "The minimum interval (in seconds) between collections from a GoBGP server.")
	flag.StringVar(&authToken, "auth.token", "anonymous", "The X-Token for accessing the exporter itself")
	flag.BoolVar(&isShowMetrics, "metrics", false, "Display available metrics")
	flag.BoolVar(&isShowVersion, "version", false, "version information")
	flag.StringVar(&logLevel, "log.level", "info", "logging severity level")

	var usageHelp = func() {
		fmt.Fprintf(os.Stderr, "\n%s - Prometheus Exporter for GoBGP\n\n", exporter.GetExporterName())
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments]\n\n", exporter.GetExporterName())
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDocumentation: https://github.com/ovnworks/gobgp_exporter/\n\n")
	}
	flag.Usage = usageHelp
	flag.Parse()

	opts := exporter.Options{
		Address: serverAddress,
		Timeout: pollTimeout,
	}

	if err := log.Base().SetLevel(logLevel); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}

	if isShowVersion {
		fmt.Fprintf(os.Stdout, "%s %s", exporter.GetExporterName(), exporter.GetVersion())
		if exporter.GetRevision() != "" {
			fmt.Fprintf(os.Stdout, ", commit: %s\n", exporter.GetRevision())
		} else {
			fmt.Fprint(os.Stdout, "\n")
		}
		os.Exit(0)
	}

	if isShowMetrics {
		e := &exporter.RouterNode{}
		fmt.Fprintf(os.Stdout, "%s\n", e.GetMetricsTable())
		os.Exit(0)
	}

	log.Infof("Starting %s %s", exporter.GetExporterName(), exporter.GetVersionInfo())
	log.Infof("Build context %s", exporter.GetVersionBuildContext())

	e, err := exporter.NewExporter(opts)
	if err != nil {
		log.Errorf("%s failed to init properly: %s", exporter.GetExporterName(), err)
		os.Exit(1)
	}
	e.SetPollInterval(int64(pollInterval))
	if err := e.AddAuthenticationToken(authToken); err != nil {
		log.Errorf("%s failed to add authentication token: %s", exporter.GetExporterName(), err)
		os.Exit(1)
	}

	// log.Infof("GoBGP Server: %s", e.ServerAddress)
	log.Infof("Minimal scrape interval: %d seconds", e.GetPollInterval())

	http.HandleFunc(metricsPath, func(w http.ResponseWriter, r *http.Request) {
		e.Scrape(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		e.Summary(metricsPath, w, r)
	})

	log.Infoln("Listening on", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

package main

import (
	"flag"
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	gobgpapi "github.com/osrg/gobgp/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"golang.org/x/net/context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"
)

const (
	namespace = "gobgp"
)

var (
	appVersion = "[untracked]"
	gitBranch  string
	gitCommit  string
	buildUser  string // whoami
	buildDate  string // date -u
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Is GoBGP up and responds to queries (1) or is it down (0).",
		nil, nil,
	)
	routerAS = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "asn"),
		"What is GoBGP router ID and AS number.",
		[]string{"router_id"}, nil,
	)
	routerLastConnected = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "connected_at"),
		"When was the last successful connection to GoBGP.",
		[]string{"router_id"}, nil,
	)
	routerLostConnection = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "lost_connection_at"),
		"When did the exporter lose connection to the router.",
		[]string{"router_id"}, nil,
	)
	/*
		bgpPeer = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "peer_state"),
			"What is BGP peer state: Established, Idle, etc.",
			[]string{"router_id"}, nil,
		)
		bgpPeerUptime = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "peer_uptime"),
			"How long BGP peer is up.",
			[]string{"router_id"}, nil,
		)
		bgpPeerReceivedRoutes = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "peer_received_routes"),
			"How many routes did the BGP peer sent.",
			[]string{"router_id", "address_family"}, nil,
		)
		bgpPeerAcceptedRoutes = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "peer_accepted_routes"),
			"How many routes were accepted from the BGP peer.",
			[]string{"router_id", "address_family"}, nil,
		)
		bgpRoutesTotal = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "routes_total"),
			"How many routes are in BGP address family.",
			[]string{"address_family"}, nil,
		)
	*/
)

// Exporter collects GoBGP data from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
	client         gobgpapi.GobgpApiExtendedClient
	address        string
	timeout        int
	lastConnected  int64
	lostConnection int64
	connected      bool
	routerID       string
	localAS        uint32
}

type gobgpOpts struct {
	address string
	timeout int
}

// NewExporter returns an initialized Exporter.
func NewExporter(opts gobgpOpts) (*Exporter, error) {
	e := Exporter{
		address: opts.address,
		timeout: opts.timeout,
	}
	client, err := gobgpapi.NewGobgpApiExtendedClient(opts.address, opts.timeout)
	e.client = client
	if err != nil {
		return &e, err
	}
	return &e, nil
}

// Describe describes all the metrics ever exported by the GoBGP exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- routerAS
	ch <- routerLastConnected
	ch <- routerLostConnection
}

// Reconnect closes existing connection to GoBGP, if any. Then, it
// creates a new one.
func (e *Exporter) Reconnect() error {
	if e.client.Conn != nil {
		e.client.Conn.Close()
	}
	e.connected = false
	client, err := gobgpapi.NewGobgpApiExtendedClient(e.address, 1)
	if err != nil {
		return err
	}
	e.client = client
	e.lastConnected = time.Now().Unix()
	e.connected = true
	return nil
}

func IsConnectionError(err error) bool {
	if strings.Contains(err.Error(), "connection is") {
		return true
	}
	return false
}

// Collect fetches the stats from GoBGP server and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	upValue := 0

	if e.connected {
		// What is RouterID and AS number of this GoBGP server?
		req := new(gobgpapi.GetServerRequest)
		server, err := e.client.Gobgp.GetServer(context.Background(), req)
		if err != nil {
			log.Errorf("Can't query GoBGP: %v", err)
			if IsConnectionError(err) {
				if e.connected {
					e.lostConnection = time.Now().Unix()
					e.connected = false
				}
				log.Errorf("Failed to connect to GoBGP: %v", err)
				if err := e.Reconnect(); err != nil {
					e.connected = false
					log.Errorf("Failed to reconnect to GoBGP: %v", err)
				}
			}
		} else {
			e.routerID = server.Global.GetRouterId()
			e.localAS = server.Global.GetAs()
			upValue = 1
		}
	} else {
		if err := e.Reconnect(); err != nil {
			log.Errorf("Failed to reconnect to GoBGP: %v", err)
		}
	}

	ch <- prometheus.MustNewConstMetric(
		up,
		prometheus.GaugeValue,
		float64(upValue),
	)

	ch <- prometheus.MustNewConstMetric(
		routerLastConnected,
		prometheus.CounterValue,
		float64(e.lastConnected),
		e.routerID,
	)

	ch <- prometheus.MustNewConstMetric(
		routerLostConnection,
		prometheus.CounterValue,
		float64(e.lostConnection),
		e.routerID,
	)

	if !e.connected {
		return
	}

	ch <- prometheus.MustNewConstMetric(
		routerAS,
		prometheus.GaugeValue,
		float64(e.localAS),
		e.routerID,
	)
}

func init() {
	prometheus.MustRegister(version.NewCollector("gobgp_exporter"))
}

func main() {
	var listenAddress, metricsPath, gobgpAddress string
	var gobgpTimeout int
	var isShowVersion bool
	appName := "gobgp_exporter"
	flag.StringVar(&listenAddress, "web.listen-address", ":9472", "Address to listen on for web interface and telemetry.")
	flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	opts := gobgpOpts{}
	flag.StringVar(&gobgpAddress, "gobgp.address", "127.0.0.1:50051", "gRPC API address of GoBGP server.")
	flag.IntVar(&gobgpTimeout, "gobgp.timeout", 2, "Timeout on gRPC requests to GoBGP.")
	flag.BoolVar(&isShowVersion, "version", false, "version information")
	opts.address = gobgpAddress
	opts.timeout = gobgpTimeout
	var usageHelp = func() {
		fmt.Fprintf(os.Stderr, "\n%s - Prometheus Exporter for GoBGP\n\n", appName)
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments]\n\n", appName)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDocumentation: https://github.com/greenpau/%s/\n\n", appName)
	}
	flag.Usage = usageHelp
	flag.Parse()
	version.Version = appVersion
	version.Revision = gitCommit
	version.Branch = gitBranch
	version.BuildUser = buildUser
	version.BuildDate = buildDate

	if isShowVersion {
		fmt.Fprintf(os.Stdout, "%s %s", appName, version.Version)
		if version.Revision != "" {
			fmt.Fprintf(os.Stdout, ", commit: %s\n", version.Revision)
		} else {
			fmt.Fprint(os.Stdout, "\n")
		}
		os.Exit(0)
	}

	log.Infoln("Starting gobgp_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	exporter, err := NewExporter(opts)
	if err != nil {
		log.Errorf("gobgp_exporter failed to init properly: %s", err)
		exporter.connected = false
	} else {
		exporter.lastConnected = time.Now().Unix()
		exporter.connected = true
	}
	prometheus.MustRegister(exporter)

	http.Handle(metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>GoBGP Exporter</title></head>
             <body>
             <h1>GoBGP Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.Infoln("Listening on", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

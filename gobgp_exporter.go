package main

import (
	"flag"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	gobgpapi "github.com/osrg/gobgp/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"golang.org/x/net/context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"sync"
	"sync/atomic"
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
	routerQueryErrors = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "failed_query_count"),
		"The number of failed queries to GoBGP router",
		[]string{"router_id"}, nil,
	)
	routerRibDestinations = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "route_count"),
		"The number of routes on per address family and resource type basis",
		[]string{"router_id", "resource_type", "address_family"}, nil,
	)
	routerNextPoll = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "next_poll"),
		"The timestamp of the next potential poll of GoBGP server",
		[]string{"router_id"}, nil,
	)
	routerPeers = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_count"),
		"The number of BGP peers",
		[]string{"router_id"}, nil,
	)
	routerPeer = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_up"),
		"Is the peer up and in established state (1) or it is not (0).",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	routerPeerAsn = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_asn"),
		"What is the AS number of the peer",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	routerPeerLocalAsn = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_local_asn"),
		"What is the AS number presented to the peer by this router.",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	routerPeerAdminState = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_admin_state"),
		"Is the peer configured for being Up (0), Down (1), or PFX_CT (2)",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	routerPeerSessionState = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_session_state"),
		"What is the state of BGP session to the peer: unknown (0), idle (1), connect (2), active (3), opensent (4), openconfirm (5), established (6)",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	bgpPeerReceivedRoutes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_received_route_count"),
		"How many routes did the BGP peer sent to this router (limited to IPv4).",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	bgpPeerAcceptedRoutes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_accepted_route_count"),
		"How many routes were accepted from the routes received from this BGP peer (limited to IPv4)",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	bgpPeerAdvertisedRoutes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_advertised_route_count"),
		"How many routes were advertised to this BGP peer (limited to IPv4).",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	bgpPeerOutQueue = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_out_queue_count"),
		"PeerState.OutQ",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	bgpPeerFlops = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_flop_count"),
		"PeerState.Flops",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	bgpPeerSendCommunityFlag = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_send_community"),
		"PeerState.SendCommunity",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	bgpPeerRemovePrivateAsFlag = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_remove_private_as"),
		"PeerState.RemovePrivateAs",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	bgpPeerPasswodSetFlag = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_password_set"),
		"Whether the GoBGP peer has been configured (1) for authentication or not (0)",
		[]string{"router_id", "peer_router_id"}, nil,
	)
	bgpPeerType = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "peer_type"),
		"PeerState.PeerType",
		[]string{"router_id", "peer_router_id"}, nil,
	)
)

// Exporter collects GoBGP data from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
	sync.RWMutex
	client          gobgpapi.GobgpApiExtendedClient
	address         string
	timeout         int
	pollInterval    int64
	lastConnected   int64
	lostConnection  int64
	errors          int64
	errorsLocker    sync.RWMutex
	connected       bool
	routerID        string
	localAS         uint32
	resourceTypes   map[string]bool
	addressFamilies map[string]bool
	lastCollected   int64
	metrics         []prometheus.Metric
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
	e.resourceTypes = make(map[string]bool)
	e.addressFamilies = make(map[string]bool)
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
	ch <- routerQueryErrors
	ch <- routerRibDestinations
	ch <- routerNextPoll
	ch <- routerPeers
	ch <- routerPeer
	ch <- routerPeerAsn
	ch <- routerPeerLocalAsn
	ch <- routerPeerAdminState
	ch <- routerPeerSessionState
	ch <- bgpPeerReceivedRoutes
	ch <- bgpPeerAcceptedRoutes
	ch <- bgpPeerAdvertisedRoutes
	ch <- bgpPeerOutQueue
	ch <- bgpPeerFlops
	ch <- bgpPeerSendCommunityFlag
	ch <- bgpPeerRemovePrivateAsFlag
	ch <- bgpPeerPasswodSetFlag
	ch <- bgpPeerType
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

// IsConnectionError checks whether it is connectivity issue.
func IsConnectionError(err error) bool {
	if strings.Contains(err.Error(), "connection is") {
		return true
	}
	return false
}

// IncrementErrorCounter increases the counter of failed queries
// to GoBGP server.
func (e *Exporter) IncrementErrorCounter() {
	e.errorsLocker.Lock()
	defer e.errorsLocker.Unlock()
	atomic.AddInt64(&e.errors, 1)
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.GatherMetrics()
	e.RLock()
	defer e.RUnlock()
	if len(e.metrics) == 0 {
		ch <- prometheus.MustNewConstMetric(
			up,
			prometheus.GaugeValue,
			0,
		)
		ch <- prometheus.MustNewConstMetric(
			routerQueryErrors,
			prometheus.CounterValue,
			float64(e.errors),
			e.routerID,
		)
		ch <- prometheus.MustNewConstMetric(
			routerNextPoll,
			prometheus.CounterValue,
			float64(e.lastCollected),
			e.routerID,
		)
		return
	}
	for _, m := range e.metrics {
		ch <- m
	}
}

// GatherMetrics collect data from GoBGP server and stores them
// as Prometheus metrics.
func (e *Exporter) GatherMetrics() {
	if time.Now().Unix() < e.lastCollected {
		return
	}
	e.Lock()
	defer e.Unlock()
	if len(e.metrics) > 0 {
		e.metrics = e.metrics[:0]
	}
	upValue := 0
	if e.connected {
		// What is RouterID and AS number of this GoBGP server?
		req := new(gobgpapi.GetServerRequest)
		server, err := e.client.Gobgp.GetServer(context.Background(), req)
		if err != nil {
			e.IncrementErrorCounter()
			log.Errorf("Can't query GoBGP: %v", err)
			if IsConnectionError(err) {
				if e.connected {
					e.lostConnection = time.Now().Unix()
					e.connected = false
				}
				log.Errorf("Failed to connect to GoBGP: %v", err)
				if err := e.Reconnect(); err != nil {
					e.IncrementErrorCounter()
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
			e.IncrementErrorCounter()
			log.Errorf("Failed to reconnect to GoBGP: %v", err)
		}
	}

	e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
		up,
		prometheus.GaugeValue,
		float64(upValue),
	))

	e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
		routerLastConnected,
		prometheus.CounterValue,
		float64(e.lastConnected),
		e.routerID,
	))

	e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
		routerLostConnection,
		prometheus.CounterValue,
		float64(e.lostConnection),
		e.routerID,
	))

	if e.connected {
		e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
			routerAS,
			prometheus.GaugeValue,
			float64(e.localAS),
			e.routerID,
		))
	}

	if e.connected {
		for _, resourceTypeName := range gobgpapi.Resource_name {
			for _, addressFamilyName := range gobgpapi.Family_name {
				if !e.connected {
					continue
				}
				if _, exists := e.resourceTypes[resourceTypeName]; !exists {
					continue
				}
				if _, exists := e.addressFamilies[addressFamilyName]; !exists {
					continue
				}
				var resourceType gobgpapi.Resource
				switch resourceTypeName {
				case "GLOBAL":
					resourceType = gobgpapi.Resource_GLOBAL
				case "LOCAL":
					resourceType = gobgpapi.Resource_LOCAL
				default:
					continue
				}
				ribRequest := new(gobgpapi.GetRibRequest)
				ribRequest.Table = &gobgpapi.Table{
					Type:   resourceType,
					Family: uint32(gobgpapi.Family_value[addressFamilyName]),
				}
				rib, err := e.client.Gobgp.GetRib(context.Background(), ribRequest)
				if err != nil {
					log.Errorf("GoBGP query failed for resource type %s for %s address family: %s", resourceTypeName, addressFamilyName, err)
					e.IncrementErrorCounter()
					continue
				}
				log.Infof("GoBGP RIB size for %s/%s: %d", resourceTypeName, addressFamilyName, len(rib.Table.Destinations))
				//log.Debugf("GoBGP RIB size for %s/%s: %d", resourceTypeName, addressFamilyName, len(rib.Table.Destinations))
				//spew.Dump(len(rib.Destinations))
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					routerRibDestinations,
					prometheus.GaugeValue,
					float64(len(rib.Table.Destinations)),
					e.routerID,
					strings.ToLower(resourceTypeName),
					strings.ToLower(addressFamilyName),
				))
			}
		}
	}

	if e.connected {
		peerRequest := new(gobgpapi.GetNeighborRequest)
		peerRequest.EnableAdvertised = false
		peerRequest.Address = ""
		peerResponse, err := e.client.Gobgp.GetNeighbor(context.Background(), peerRequest)
		if err != nil {
			log.Errorf("GoBGP query for peers failed: %s", err)
			e.IncrementErrorCounter()
		} else {
			e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
				routerPeers,
				prometheus.GaugeValue,
				float64(len(peerResponse.Peers)),
				e.routerID,
			))
			for _, p := range peerResponse.Peers {
				peerRouterID := p.Info.NeighborAddress
				// Peer Up/Down
				if strings.HasSuffix(p.Info.BgpState, "stablished") {
					e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
						routerPeer,
						prometheus.GaugeValue,
						1,
						e.routerID,
						peerRouterID,
					))
				} else {
					e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
						routerPeer,
						prometheus.GaugeValue,
						0,
						e.routerID,
						peerRouterID,
					))
				}
				// Peer ASN
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					routerPeerAsn,
					prometheus.GaugeValue,
					float64(p.Info.PeerAs),
					e.routerID,
					peerRouterID,
				))
				// Peer Admin State: Up (0), Down (1), PFX_CT (2)
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					routerPeerAdminState,
					prometheus.GaugeValue,
					float64(p.Info.AdminState),
					e.routerID,
					peerRouterID,
				))
				// Peer Session State: idle (0), connect (1), active (2), opensent (3)
				// openconfirm (4), established (5).
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					routerPeerSessionState,
					prometheus.GaugeValue,
					float64(p.Info.SessionState),
					e.routerID,
					peerRouterID,
				))
				// Local AS advertised to the peer
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					routerPeerLocalAsn,
					prometheus.GaugeValue,
					float64(p.Info.LocalAs),
					e.routerID,
					peerRouterID,
				))
				// The number of received routes from the peer
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					bgpPeerReceivedRoutes,
					prometheus.GaugeValue,
					float64(p.Info.Received),
					e.routerID,
					peerRouterID,
				))
				// The number of accepted routes from the peer
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					bgpPeerAcceptedRoutes,
					prometheus.GaugeValue,
					float64(p.Info.Accepted),
					e.routerID,
					peerRouterID,
				))
				// The number of advertised routes to the peer
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					bgpPeerAdvertisedRoutes,
					prometheus.GaugeValue,
					float64(p.Info.Advertised),
					e.routerID,
					peerRouterID,
				))
				// TODO: PeerState.OutQ
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					bgpPeerOutQueue,
					prometheus.GaugeValue,
					float64(p.Info.OutQ),
					e.routerID,
					peerRouterID,
				))
				// TODO: PeerState.Flops
				// TODO: Is it a gauge or counter?
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					bgpPeerFlops,
					prometheus.GaugeValue,
					float64(p.Info.Flops),
					e.routerID,
					peerRouterID,
				))
				// TODO: PeerState.SendCommunity
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					bgpPeerSendCommunityFlag,
					prometheus.GaugeValue,
					float64(p.Info.SendCommunity),
					e.routerID,
					peerRouterID,
				))
				// TODO: PeerState.RemovePrivateAs
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					bgpPeerRemovePrivateAsFlag,
					prometheus.GaugeValue,
					float64(p.Info.RemovePrivateAs),
					e.routerID,
					peerRouterID,
				))
				// TODO: PeerState.AuthPassword
				passwordSetFlag := 0
				if p.Info.AuthPassword != "" {
					passwordSetFlag = 1
				}
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					bgpPeerPasswodSetFlag,
					prometheus.GaugeValue,
					float64(passwordSetFlag),
					e.routerID,
					peerRouterID,
				))
				// TODO: PeerState.PeerType
				e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
					bgpPeerType,
					prometheus.GaugeValue,
					float64(p.Info.PeerType),
					e.routerID,
					peerRouterID,
				))

			}
		}
	}

	e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
		routerQueryErrors,
		prometheus.CounterValue,
		float64(e.errors),
		e.routerID,
	))

	e.metrics = append(e.metrics, prometheus.MustNewConstMetric(
		routerNextPoll,
		prometheus.CounterValue,
		float64(e.lastCollected),
		e.routerID,
	))
	e.lastCollected = time.Now().Add(time.Duration(e.pollInterval) * time.Second).Unix()
}

func init() {
	prometheus.MustRegister(version.NewCollector("gobgp_exporter"))
}

func main() {
	var listenAddress, metricsPath, gobgpAddress string
	var gobgpTimeout, gobgpPollInterval int
	var isShowVersion bool
	appName := "gobgp_exporter"
	flag.StringVar(&listenAddress, "web.listen-address", ":9472", "Address to listen on for web interface and telemetry.")
	flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	opts := gobgpOpts{}
	flag.StringVar(&gobgpAddress, "gobgp.address", "127.0.0.1:50051", "gRPC API address of GoBGP server.")
	flag.IntVar(&gobgpTimeout, "gobgp.timeout", 2, "Timeout on gRPC requests to GoBGP.")
	flag.IntVar(&gobgpPollInterval, "gobgp.poll-interval", 15, "The minimum interval (in seconds) between collections from GoBGP server.")
	flag.BoolVar(&isShowVersion, "version", false, "version information")
	var usageHelp = func() {
		fmt.Fprintf(os.Stderr, "\n%s - Prometheus Exporter for GoBGP\n\n", appName)
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments]\n\n", appName)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDocumentation: https://github.com/greenpau/%s/\n\n", appName)
	}
	flag.Usage = usageHelp
	flag.Parse()
	opts.address = gobgpAddress
	opts.timeout = gobgpTimeout
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
	log.Infof("GoBGP server: %s", opts.address)

	exporter, err := NewExporter(opts)
	if err != nil {
		log.Errorf("gobgp_exporter failed to init properly: %s", err)
		exporter.connected = false
	} else {
		exporter.lastConnected = time.Now().Unix()
		exporter.connected = true
	}
	exporter.pollInterval = int64(gobgpPollInterval)
	exporter.resourceTypes["LOCAL"] = true
	exporter.resourceTypes["GLOBAL"] = true
	exporter.addressFamilies["IPv4"] = true
	exporter.addressFamilies["EVPN"] = true
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

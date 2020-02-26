// Copyright 2018 Paul Greenberg (greenpau@outlook.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exporter

import (
	"fmt"
	gobgpapi "github.com/osrg/gobgp/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type credential struct {
	Username string
	Password string
	Failed   bool
}

// RouterNode is an instance of a GoBGP router.
type RouterNode struct {
	sync.RWMutex
	client               gobgpapi.GobgpApiClient
	address              string
	routerID             string
	localAS              uint32
	resourceTypes        map[string]bool
	addressFamilies      map[string]bool
	port                 int
	result               string
	module               string
	timestamp            string
	pollInterval         int64
	timeout              int
	errors               int64
	errorsLocker         sync.RWMutex
	nextCollectionTicker int64
	metrics              []prometheus.Metric
	lastConnected        int64
	connected            bool
}

// NewRouterNode creates an instance of RouterNode.
func NewRouterNode(addr string, timeout int) (*RouterNode, error) {
	if err := validAddress(addr); err != nil {
		return nil, err
	}
	n := &RouterNode{
		result:               "unknown",
		timestamp:            "unknown",
		nextCollectionTicker: 0,
		errors:               0,
		address:              addr,
	}
	n.resourceTypes = make(map[string]bool)
	n.addressFamilies = make(map[string]bool)
	n.resourceTypes["LOCAL"] = true
	n.resourceTypes["GLOBAL"] = true
	n.addressFamilies["IPv4"] = true
	n.addressFamilies["EVPN"] = true

	grpcOpts := []grpc.DialOption{grpc.WithBlock()}
	grpcOpts = append(grpcOpts, grpc.WithInsecure())

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpcOpts...)
	if err != nil {
		n.IncrementErrorCounter()
		return n, err
	}

	n.client = gobgpapi.NewGobgpApiClient(conn)
	return n, nil
}

func validAddress(s string) error {
	if s == "" {
		return fmt.Errorf("empty address")
	}
	arr := strings.Split(s, ":")
	if len(arr) != 2 {
		return fmt.Errorf("invalid address %s, expected 'ipaddress:port'", s)
	}
	if addr := net.ParseIP(arr[0]); addr == nil {
		return fmt.Errorf("invalid IP address in %s", s)
	}
	port, err := strconv.Atoi(arr[1])
	if err != nil {
		return err
	}
	if strconv.Itoa(port) != arr[1] {
		return fmt.Errorf("invalid port in %s", s)
	}
	if port < 1024 || port > 65535 {
		return fmt.Errorf("invalid port in %s, expected range 1024-65535", s)
	}
	return nil
}

// IncrementErrorCounter increases the counter of failed queries
// to a network node.
func (n *RouterNode) IncrementErrorCounter() {
	n.errorsLocker.Lock()
	defer n.errorsLocker.Unlock()
	atomic.AddInt64(&n.errors, 1)
}

// Collect implements prometheus.Collector.
func (n *RouterNode) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	log.Debugf("Calling GatherMetrics()")
	n.GatherMetrics()
	log.Debugf("Collect() calls RLock()")
	n.RLock()
	defer n.RUnlock()
	log.Debugf("Collect() successful RLock()")
	if len(n.metrics) == 0 {
		log.Debugf("Collect() no metrics found")
		ch <- prometheus.MustNewConstMetric(
			routerUp,
			prometheus.GaugeValue,
			0,
		)
		ch <- prometheus.MustNewConstMetric(
			routerErrors,
			prometheus.CounterValue,
			float64(n.errors),
		)
		ch <- prometheus.MustNewConstMetric(
			routerNextScrape,
			prometheus.CounterValue,
			float64(n.nextCollectionTicker),
		)
		ch <- prometheus.MustNewConstMetric(
			routerScrapeTime,
			prometheus.GaugeValue,
			time.Since(start).Seconds(),
		)
		return
	}
	log.Debugf("Collect() sends %d metrics to a shared channel", len(n.metrics))
	for _, m := range n.metrics {
		ch <- m
	}
}

// IsConnectionError checks whether it is connectivity issue.
func IsConnectionError(err error) bool {
	if strings.Contains(err.Error(), "connection") {
		return true
	}
	return false
}

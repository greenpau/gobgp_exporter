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
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	gobgpapi "github.com/osrg/gobgp/v3/api"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// RouterNode is an instance of a GoBGP router.
type RouterNode struct {
	sync.RWMutex
	client               gobgpapi.GobgpApiClient
	address              string
	routerID             string
	localAS              uint32
	resourceTypes        map[string]bool
	addressFamilies      map[string]bool
	result               string
	timestamp            string
	pollInterval         int64
	errors               int64
	errorsLocker         sync.RWMutex
	nextCollectionTicker int64
	metrics              []prometheus.Metric
	connected            bool
	logger               log.Logger
}

// NewRouterNode creates an instance of RouterNode.
func NewRouterNode(addr string, timeout int, tlsConfig *tls.Config, logger log.Logger) (*RouterNode, error) {
	if err := validAddress(addr, logger); err != nil {
		return nil, err
	}
	n := &RouterNode{
		result:               "unknown",
		timestamp:            "unknown",
		nextCollectionTicker: 0,
		errors:               0,
		address:              addr,
		logger:               logger,
	}
	n.resourceTypes = make(map[string]bool)
	n.addressFamilies = make(map[string]bool)
	n.resourceTypes["LOCAL"] = true
	n.resourceTypes["GLOBAL"] = true
	n.addressFamilies["IPv4"] = true
	n.addressFamilies["EVPN"] = true

	grpcOpts := []grpc.DialOption{grpc.WithBlock()}
	if tlsConfig == nil {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	}

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

func validAddress(s string, logger log.Logger) error {
	if s == "" {
		return fmt.Errorf("empty address")
	}

	host, strport, err := net.SplitHostPort(s)
	if err != nil {
		return err
	} else if host != "" {
		if addr := net.ParseIP(host); addr == nil {
			return fmt.Errorf("invalid IP address in %s", s)
		}
	} else if !strings.HasPrefix(s, "dns://") {
		return fmt.Errorf("invalid address format in %s", s)
	} else {
		// "dns://" prefix for hostname is allowed per go grpc documentation
		// see https://pkg.go.dev/google.golang.org/grpc#DialContext
		idx := strings.LastIndex(s, ":")
		host = s[0:idx]
		strport = s[idx+1:]
	}

	level.Debug(logger).Log(
		"msg", "validAddress info",
		"uri", s,
		"host", host,
		"port", strport,
	)

	port, err := strconv.Atoi(strport)
	if err != nil {
		return err
	}
	if strconv.Itoa(port) != strport {
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
	level.Debug(n.logger).Log(
		"msg", "Calling GatherMetrics()",
	)
	n.GatherMetrics()
	level.Debug(n.logger).Log(
		"msg", "Collect() calls RLock()",
	)
	n.RLock()
	defer n.RUnlock()
	level.Debug(n.logger).Log(
		"msg", "Collect() successful RLock()",
	)
	if len(n.metrics) == 0 {
		level.Debug(n.logger).Log(
			"msg", "Collect() no metrics found",
		)
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
	level.Debug(n.logger).Log(
		"msg", "Collect() sends metrics to a shared channel",
		"metric_count", len(n.metrics),
	)
	for _, m := range n.metrics {
		ch <- m
	}
}

// IsConnectionError checks whether it is connectivity issue.
func IsConnectionError(err error) bool {
	return strings.Contains(err.Error(), "connection")
}

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
	"sync"
	"time"

	"github.com/go-kit/log/level"
	gobgpapi "github.com/osrg/gobgp/v3/api"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
)

// GatherMetrics collect data from a GoBGP router and stores them
// as Prometheus metrics.
func (n *RouterNode) GatherMetrics() {
	n.Lock()
	defer n.Unlock()

	level.Debug(n.logger).Log(
		"msg", "GatherMetrics() locked",
	)

	if time.Now().Unix() < n.nextCollectionTicker {
		return
	}
	start := time.Now()
	if len(n.metrics) > 0 {
		n.metrics = n.metrics[:0]
		level.Debug(n.logger).Log(
			"msg", "GatherMetrics() cleared metrics",
		)
	}
	upValue := 1

	// What is RouterID and AS number of this GoBGP server?
	server, err := n.client.GetBgp(context.Background(), &gobgpapi.GetBgpRequest{})
	if err != nil {
		n.IncrementErrorCounter()
		level.Error(n.logger).Log(
			"msg", "failed query gobgp server",
			"error", err.Error(),
		)
		if IsConnectionError(err) {
			n.connected = false
			upValue = 0
		}
	} else {
		n.routerID = server.Global.RouterId
		n.localAS = server.Global.Asn
		level.Debug(n.logger).Log(
			"msg", "router info",
			"router_id", n.routerID,
			"local_asn", n.localAS,
		)
		n.connected = true
	}

	if n.connected {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			n.GetRibCounters()
		}()
		go func() {
			defer wg.Done()
			n.GetPeers()
		}()
		wg.Wait()

	}

	// Generic Metrics
	n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
		routerUp,
		prometheus.GaugeValue,
		float64(upValue),
	))

	n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
		routerErrors,
		prometheus.CounterValue,
		float64(n.errors),
	))
	n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
		routerNextScrape,
		prometheus.CounterValue,
		float64(n.nextCollectionTicker),
	))
	n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
		routerScrapeTime,
		prometheus.GaugeValue,
		time.Since(start).Seconds(),
	))

	// Router ID and ASN
	if n.routerID != "" {
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			routerID,
			prometheus.GaugeValue,
			1,
			n.routerID,
		))
	}
	if n.localAS > 0 {
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			routerLocalAS,
			prometheus.GaugeValue,
			float64(n.localAS),
		))
	}

	n.nextCollectionTicker = time.Now().Add(time.Duration(n.pollInterval) * time.Second).Unix()

	if upValue > 0 {
		n.result = "success"
	} else {
		n.result = "failure"
	}
	n.timestamp = time.Now().Format(time.RFC3339)

	level.Debug(n.logger).Log(
		"msg", "GatherMetrics() returns",
	)
}

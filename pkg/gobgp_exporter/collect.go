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
	//"github.com/davecgh/go-spew/spew"
	gobgpapi "github.com/osrg/gobgp/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"golang.org/x/net/context"
	"sync"
	"time"
)

// GatherMetrics collect data from a GoBGP router and stores them
// as Prometheus metrics.
func (n *RouterNode) GatherMetrics() {
	n.Lock()
	defer n.Unlock()
	log.Debugf("GatherMetrics() locked")
	if time.Now().Unix() < n.nextCollectionTicker {
		return
	}
	start := time.Now()
	if len(n.metrics) > 0 {
		n.metrics = n.metrics[:0]
		log.Debugf("GatherMetrics() cleared metrics")
	}
	upValue := 1

	if n.connected {
		// What is RouterID and AS number of this GoBGP server?
		req := new(gobgpapi.GetServerRequest)
		server, err := n.client.Gobgp.GetServer(context.Background(), req)
		if err != nil {
			n.IncrementErrorCounter()
			log.Errorf("Can't query GoBGP: %v", err)
			if IsConnectionError(err) {
				n.connected = false
				if err := n.Reconnect(); err != nil {
					n.IncrementErrorCounter()
					n.connected = false
					log.Errorf("Failed to reconnect to GoBGP: %v", err)
					upValue = 0
				}
			}
		}
		if n.connected {
			n.routerID = server.Global.GetRouterId()
			n.localAS = server.Global.GetAs()
		}
	} else {
		if err := n.Reconnect(); err != nil {
			n.IncrementErrorCounter()
			log.Errorf("Failed to reconnect to GoBGP: %v", err)
			upValue = 0
		} else {
			log.Debugf("Router ID: '%s', ASN: %d", n.routerID, n.localAS)
		}
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

	log.Debugf("GatherMetrics() returns")
	return
}

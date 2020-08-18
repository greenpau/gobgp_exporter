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
	"github.com/prometheus/client_golang/prometheus"
)

var (
	routerUp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "router", "up"),
		"Is GoBGP up and responds to queries (1) or is it down (0).",
		nil, nil,
	)
	routerID = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "router", "id"),
		"What is GoBGP router ID.",
		[]string{"id"}, nil,
	)
	routerLocalAS = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "router", "asn"),
		"What is GoBGP AS number.",
		nil, nil,
	)
	routerErrors = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "router", "failed_req_count"),
		"The number of failed requests to GoBGP router.",
		nil, nil,
	)
	routerNextScrape = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "router", "next_poll"),
		"The timestamp of the next potential scrape of the router.",
		nil, nil,
	)
	routerScrapeTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "router", "scrape_time"),
		"The amount of time it took to scrape the router.",
		nil, nil,
	)
	routerPeers = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "count"),
		"The number of BGP peers",
		nil, nil,
	)
)

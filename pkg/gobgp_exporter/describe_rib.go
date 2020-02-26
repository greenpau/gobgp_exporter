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
	routerRibTotalDestinationCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "route", "total_destination_count"),
		"The number of routes on per address family and route table basis",
		[]string{"route_table", "address_family", "vrf_name"}, nil,
	)

	routerRibTotalPathCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "route", "total_path_count"),
		"The number of available paths to destinations on per address family and route table basis",
		[]string{"route_table", "address_family", "vrf_name"}, nil,
	)

	routerRibAcceptedPathCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "route", "accepted_path_count"),
		"The number of accepted paths to destinations on per address family and route table basis",
		[]string{"route_table", "address_family", "vrf_name"}, nil,
	)
)

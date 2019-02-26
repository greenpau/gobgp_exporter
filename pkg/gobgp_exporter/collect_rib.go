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
	gobgpapi "github.com/osrg/gobgp/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"golang.org/x/net/context"
	"strings"
)

// GetRibCounters collects BGP routing information base (RIB) related metrics.
func (n *RouterNode) GetRibCounters() {
	if n.connected == false {
		return
	}
	for _, resourceTypeName := range gobgpapi.Resource_name {
		for _, addressFamilyName := range gobgpapi.Family_name {
			if !n.connected {
				continue
			}
			if _, exists := n.resourceTypes[resourceTypeName]; !exists {
				continue
			}
			if _, exists := n.addressFamilies[addressFamilyName]; !exists {
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
			rib, err := n.client.Gobgp.GetRib(context.Background(), ribRequest)
			if err != nil {
				log.Errorf("GoBGP query failed for resource type %s for %s address family: %s", resourceTypeName, addressFamilyName, err)
				n.IncrementErrorCounter()
				continue
			}
			log.Debugf("GoBGP RIB size for %s/%s: %d", resourceTypeName, addressFamilyName, len(rib.Table.Destinations))
			//spew.Dump(len(rib.Destinations))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				routerRibDestinations,
				prometheus.GaugeValue,
				float64(len(rib.Table.Destinations)),
				strings.ToLower(resourceTypeName),
				strings.ToLower(addressFamilyName),
			))
		}
	}
	return
}

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
	gobgpapi "github.com/osrg/gobgp/v3/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"golang.org/x/net/context"
	"strings"
)

var addressFamilies = map[string]*gobgpapi.Family{
	"ipv4": {
		Afi:  gobgpapi.Family_AFI_IP,
		Safi: gobgpapi.Family_SAFI_UNICAST,
	},
	"ipv6": {
		Afi:  gobgpapi.Family_AFI_IP6,
		Safi: gobgpapi.Family_SAFI_UNICAST,
	},
	"ipv4_vpn": {
		Afi:  gobgpapi.Family_AFI_IP,
		Safi: gobgpapi.Family_SAFI_MPLS_VPN,
	},
	"ipv6_vpn": {
		Afi:  gobgpapi.Family_AFI_IP6,
		Safi: gobgpapi.Family_SAFI_MPLS_VPN,
	},
	"ipv4_mpls": {
		Afi:  gobgpapi.Family_AFI_IP,
		Safi: gobgpapi.Family_SAFI_MPLS_LABEL,
	},
	"ipv6_mpls": {
		Afi:  gobgpapi.Family_AFI_IP6,
		Safi: gobgpapi.Family_SAFI_MPLS_LABEL,
	},
	"evpn": {
		Afi:  gobgpapi.Family_AFI_L2VPN,
		Safi: gobgpapi.Family_SAFI_EVPN,
	},
	"ipv4_encap": {
		Afi:  gobgpapi.Family_AFI_IP,
		Safi: gobgpapi.Family_SAFI_ENCAPSULATION,
	},
	"ipv6_encap": {
		Afi:  gobgpapi.Family_AFI_IP6,
		Safi: gobgpapi.Family_SAFI_ENCAPSULATION,
	},
	"ipv4_flowspec": {
		Afi:  gobgpapi.Family_AFI_IP,
		Safi: gobgpapi.Family_SAFI_FLOW_SPEC_UNICAST,
	},
	"ipv6_flowspec": {
		Afi:  gobgpapi.Family_AFI_IP6,
		Safi: gobgpapi.Family_SAFI_FLOW_SPEC_UNICAST,
	},
	"ipv4_vpn_flowspec": {
		Afi:  gobgpapi.Family_AFI_IP,
		Safi: gobgpapi.Family_SAFI_FLOW_SPEC_VPN,
	},
	"ipv6_vpn_flowspec": {
		Afi:  gobgpapi.Family_AFI_IP6,
		Safi: gobgpapi.Family_SAFI_FLOW_SPEC_VPN,
	},
	"l2_vpn_flowspec": {
		Afi:  gobgpapi.Family_AFI_L2VPN,
		Safi: gobgpapi.Family_SAFI_FLOW_SPEC_VPN,
	},
}

// GetRibCounters collects BGP routing information base (RIB) related metrics.
func (n *RouterNode) GetRibCounters() {
	var tableType gobgpapi.TableType
	for tableTypeName := range gobgpapi.TableType_value {
		switch tableTypeName {
		case "GLOBAL":
			tableType = gobgpapi.TableType_GLOBAL
		case "LOCAL":
			tableType = gobgpapi.TableType_LOCAL
		case "ADJ_IN":
			// tableType = gobgpapi.TableType_ADJ_IN
			continue
		case "ADJ_OUT":
			// tableType = gobgpapi.TableType_ADJ_OUT
			continue
		case "VRF":
			//tableType = gobgpapi.TableType_VRF
			continue
		default:
			log.Warnf("Unsupported GoBGP route table type: %s", tableTypeName)
			continue
		}

		for addressFamilyName, addressFamily := range addressFamilies {
			serverResponse, err := n.client.GetTable(context.Background(), &gobgpapi.GetTableRequest{
				TableType: tableType,
				Family:    addressFamily,
				Name:      "",
			})

			if err != nil {
				log.Errorf("GoBGP query for route table %s/%s failed: %s", tableTypeName, addressFamilyName, err)
				n.IncrementErrorCounter()
				continue
			}

			if serverResponse == nil {
				log.Warnf("GoBGP route table %s/%s response is empty", tableTypeName, addressFamilyName)
				continue
			}

			log.Debugf("GoBGP route table %s/%s: %v", tableTypeName, addressFamilyName, serverResponse)

			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				routerRibTotalDestinationCount,
				prometheus.GaugeValue,
				float64(serverResponse.GetNumDestination()),
				strings.ToLower(tableTypeName),
				strings.ToLower(addressFamilyName),
				"default",
			))

			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				routerRibTotalPathCount,
				prometheus.GaugeValue,
				float64(serverResponse.GetNumPath()),
				strings.ToLower(tableTypeName),
				strings.ToLower(addressFamilyName),
				"default",
			))

			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				routerRibAcceptedPathCount,
				prometheus.GaugeValue,
				float64(serverResponse.GetNumAccepted()),
				strings.ToLower(tableTypeName),
				strings.ToLower(addressFamilyName),
				"default",
			))

		}

	}

}

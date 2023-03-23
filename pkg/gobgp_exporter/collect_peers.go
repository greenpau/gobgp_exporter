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
	"io"

	"github.com/go-kit/log/level"
	gobgpapi "github.com/osrg/gobgp/v3/api"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
)

// GetPeers collects information about BGP peers.
func (n *RouterNode) GetPeers() {

	serverResponse, err := n.client.ListPeer(context.Background(), &gobgpapi.ListPeerRequest{})
	if err != nil {
		level.Error(n.logger).Log(
			"msg", "GoBGP query for peers failed",
			"error", err.Error(),
		)
		n.IncrementErrorCounter()
		return
	}

	peers := make([]*gobgpapi.Peer, 0, 1024)
	for {
		r, err := serverResponse.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			level.Error(n.logger).Log(
				"msg", "GoBGP get neighbor response parsing failed",
				"error", err.Error(),
			)
			n.IncrementErrorCounter()
			return
		}
		peers = append(peers, r.Peer)
	}

	if err != nil {
		level.Error(n.logger).Log(
			"msg", "GoBGP query for peers failed",
			"error", err.Error(),
		)
		n.IncrementErrorCounter()
		return
	}

	n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
		routerPeers,
		prometheus.GaugeValue,
		float64(len(peers)),
	))

	for _, p := range peers {
		peerState := p.GetState()
		peerRouterID := peerState.GetNeighborAddress()

		// Peer Up/Down
		if peerState.GetRouterId() != "" {
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				routerPeer,
				prometheus.GaugeValue,
				1,
				peerRouterID,
			))
		} else {
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				routerPeer,
				prometheus.GaugeValue,
				0,
				peerRouterID,
			))
		}
		// Peer ASN
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			routerPeerAsn,
			prometheus.GaugeValue,
			float64(peerState.GetPeerAsn()),
			peerRouterID,
		))
		// Peer Admin State: Up (0), Down (1), PFX_CT (2)
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			routerPeerAdminState,
			prometheus.GaugeValue,
			float64(peerState.GetAdminState()),
			peerRouterID,
		))
		// Peer Session State: idle (0), connect (1), active (2), opensent (3)
		// openconfirm (4), established (5).
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			routerPeerSessionState,
			prometheus.GaugeValue,
			float64(peerState.GetSessionState()),
			peerRouterID,
		))
		// Local AS advertised to the peer
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			routerPeerLocalAsn,
			prometheus.GaugeValue,
			float64(peerState.GetLocalAsn()),
			peerRouterID,
		))

		peerMessages := peerState.GetMessages()
		if peerMessages != nil {
			// The number of received messages from the peer
			var peerReceivedTotalMessagesCount uint64 = 0
			var peerReceivedNotificationMessagesCount uint64 = 0
			var peerReceivedUpdateMessagesCount uint64 = 0
			var peerReceivedOpenMessagesCount uint64 = 0
			var peerReceivedKeepaliveMessagesCount uint64 = 0
			var peerReceivedRefreshMessagesCount uint64 = 0
			var peerReceivedWithdrawUpdateMessagesCount uint64 = 0
			var peerReceivedWithdrawPrefixMessagesCount uint64 = 0

			peerReceivedMessages := peerMessages.GetReceived()
			if peerReceivedMessages != nil {
				peerReceivedTotalMessagesCount = peerReceivedMessages.Total
				peerReceivedNotificationMessagesCount = peerReceivedMessages.Notification
				peerReceivedUpdateMessagesCount = peerReceivedMessages.Update
				peerReceivedOpenMessagesCount = peerReceivedMessages.Open
				peerReceivedKeepaliveMessagesCount = peerReceivedMessages.Keepalive
				peerReceivedRefreshMessagesCount = peerReceivedMessages.Refresh
				peerReceivedWithdrawUpdateMessagesCount = peerReceivedMessages.WithdrawUpdate
				peerReceivedWithdrawPrefixMessagesCount = peerReceivedMessages.WithdrawPrefix
			}

			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerReceivedTotalMessagesCount,
				prometheus.GaugeValue,
				float64(peerReceivedTotalMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerReceivedNotificationMessagesCount,
				prometheus.GaugeValue,
				float64(peerReceivedNotificationMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerReceivedUpdateMessagesCount,
				prometheus.GaugeValue,
				float64(peerReceivedUpdateMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerReceivedOpenMessagesCount,
				prometheus.GaugeValue,
				float64(peerReceivedOpenMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerReceivedKeepaliveMessagesCount,
				prometheus.GaugeValue,
				float64(peerReceivedKeepaliveMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerReceivedRefreshMessagesCount,
				prometheus.GaugeValue,
				float64(peerReceivedRefreshMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerReceivedWithdrawUpdateMessagesCount,
				prometheus.GaugeValue,
				float64(peerReceivedWithdrawUpdateMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerReceivedWithdrawPrefixMessagesCount,
				prometheus.GaugeValue,
				float64(peerReceivedWithdrawPrefixMessagesCount),
				peerRouterID,
			))

			// The number of messages sent to the peer
			var peerSentTotalMessagesCount uint64 = 0
			var peerSentNotificationMessagesCount uint64 = 0
			var peerSentUpdateMessagesCount uint64 = 0
			var peerSentOpenMessagesCount uint64 = 0
			var peerSentKeepaliveMessagesCount uint64 = 0
			var peerSentRefreshMessagesCount uint64 = 0
			var peerSentWithdrawUpdateMessagesCount uint64 = 0
			var peerSentWithdrawPrefixMessagesCount uint64 = 0

			peerSentMessages := peerMessages.GetSent()
			if peerSentMessages != nil {
				peerSentTotalMessagesCount = peerSentMessages.Total
				peerSentNotificationMessagesCount = peerSentMessages.Notification
				peerSentUpdateMessagesCount = peerSentMessages.Update
				peerSentOpenMessagesCount = peerSentMessages.Open
				peerSentKeepaliveMessagesCount = peerSentMessages.Keepalive
				peerSentRefreshMessagesCount = peerSentMessages.Refresh
				peerSentWithdrawUpdateMessagesCount = peerSentMessages.WithdrawUpdate
				peerSentWithdrawPrefixMessagesCount = peerSentMessages.WithdrawPrefix
			}

			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerSentTotalMessagesCount,
				prometheus.GaugeValue,
				float64(peerSentTotalMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerSentNotificationMessagesCount,
				prometheus.GaugeValue,
				float64(peerSentNotificationMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerSentUpdateMessagesCount,
				prometheus.GaugeValue,
				float64(peerSentUpdateMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerSentOpenMessagesCount,
				prometheus.GaugeValue,
				float64(peerSentOpenMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerSentKeepaliveMessagesCount,
				prometheus.GaugeValue,
				float64(peerSentKeepaliveMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerSentRefreshMessagesCount,
				prometheus.GaugeValue,
				float64(peerSentRefreshMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerSentWithdrawUpdateMessagesCount,
				prometheus.GaugeValue,
				float64(peerSentWithdrawUpdateMessagesCount),
				peerRouterID,
			))
			n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
				bgpPeerSentWithdrawPrefixMessagesCount,
				prometheus.GaugeValue,
				float64(peerSentWithdrawPrefixMessagesCount),
				peerRouterID,
			))

		}

		// The outbound queue message size
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			bgpPeerOutQueue,
			prometheus.GaugeValue,
			float64(peerState.GetOutQ()),
			peerRouterID,
		))
		// The number of neighbor flops
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			bgpPeerFlops,
			prometheus.GaugeValue,
			float64(peerState.GetFlops()),
			peerRouterID,
		))
		// Whether BGP community is being sent
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			bgpPeerSendCommunityFlag,
			prometheus.GaugeValue,
			float64(peerState.GetSendCommunity()),
			peerRouterID,
		))
		// Whether BGP Private AS is being removed
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			bgpPeerRemovePrivateAsFlag,
			prometheus.GaugeValue,
			float64(peerState.GetRemovePrivate()),
			peerRouterID,
		))
		// Whether authentication password is being set (1) or not (0)
		passwordSetFlag := 0
		if peerState.GetAuthPassword() != "" {
			passwordSetFlag = 1
		}
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			bgpPeerPasswodSetFlag,
			prometheus.GaugeValue,
			float64(passwordSetFlag),
			peerRouterID,
		))
		// Peer Type
		n.metrics = append(n.metrics, prometheus.MustNewConstMetric(
			bgpPeerType,
			prometheus.GaugeValue,
			float64(peerState.GetType()),
			peerRouterID,
		))

	}
}

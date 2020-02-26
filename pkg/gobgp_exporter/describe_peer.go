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
	routerPeer = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "up"),
		"Is the peer up and in established state (1) or it is not (0).",
		[]string{"name"}, nil,
	)
	routerPeerAsn = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "asn"),
		"What is the AS number of the peer",
		[]string{"name"}, nil,
	)
	routerPeerLocalAsn = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "local_asn"),
		"What is the AS number presented to the peer by this router.",
		[]string{"name"}, nil,
	)
	routerPeerAdminState = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "admin_state"),
		"Is the peer configured for being Up (0), Down (1), or PFX_CT (2)",
		[]string{"name"}, nil,
	)
	routerPeerSessionState = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "session_state"),
		"What is the state of BGP session to the peer - unknown (0), idle (1), connect (2), active (3), opensent (4), openconfirm (5), established (6)",
		[]string{"name"}, nil,
	)

	bgpPeerReceivedTotalMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "received_message_total_count"),
		"The total number of messages the BGP peer sent to this router (limited to IPv4).",
		[]string{"name"}, nil,
	)

	bgpPeerReceivedNotificationMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "received_notification_message_count"),
		"How many Notification messages did the BGP peer sent to this router (limited to IPv4).",
		[]string{"name"}, nil,
	)

	bgpPeerReceivedUpdateMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "received_update_message_count"),
		"How many Update messages did the BGP peer sent to this router (limited to IPv4).",
		[]string{"name"}, nil,
	)
	bgpPeerReceivedOpenMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "received_open_message_count"),
		"How many Open messages did the BGP peer sent to this router (limited to IPv4).",
		[]string{"name"}, nil,
	)
	bgpPeerReceivedKeepaliveMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "received_keepalive_message_count"),
		"How many messages did the BGP peer sent to this router (limited to IPv4).",
		[]string{"name"}, nil,
	)
	bgpPeerReceivedRefreshMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "received_refresh_message_count"),
		"How many Refresh messages did the BGP peer sent to this router (limited to IPv4).",
		[]string{"name"}, nil,
	)
	bgpPeerReceivedWithdrawUpdateMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "received_withdraw_update_message_count"),
		"How many WithdrawUpdate messages did the BGP peer sent to this router (limited to IPv4).",
		[]string{"name"}, nil,
	)
	bgpPeerReceivedWithdrawPrefixMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "received_withdraw_prefix_message_count"),
		"How many messages did the BGP peer sent to this router (limited to IPv4).",
		[]string{"name"}, nil,
	)

	bgpPeerSentTotalMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "sent_message_total_count"),
		"The total number of messages this router sent to this BGP peer (limited to IPv4).",
		[]string{"name"}, nil,
	)

	bgpPeerSentNotificationMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "sent_notification_message_count"),
		"How many Notification messages did this router sent to this BGP peer (limited to IPv4).",
		[]string{"name"}, nil,
	)

	bgpPeerSentUpdateMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "sent_update_message_count"),
		"How many Update messages did this router sent to this BGP peer (limited to IPv4).",
		[]string{"name"}, nil,
	)
	bgpPeerSentOpenMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "sent_open_message_count"),
		"How many Open messages did this router sent to this BGP peer (limited to IPv4).",
		[]string{"name"}, nil,
	)
	bgpPeerSentKeepaliveMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "sent_keepalive_message_count"),
		"How many messages did this router sent to this BGP peer (limited to IPv4).",
		[]string{"name"}, nil,
	)
	bgpPeerSentRefreshMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "sent_refresh_message_count"),
		"How many Refresh messages did this router sent to this BGP peer (limited to IPv4).",
		[]string{"name"}, nil,
	)
	bgpPeerSentWithdrawUpdateMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "sent_withdraw_update_message_count"),
		"How many WithdrawUpdate messages did this router sent to this BGP peer (limited to IPv4).",
		[]string{"name"}, nil,
	)
	bgpPeerSentWithdrawPrefixMessagesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "sent_withdraw_prefix_message_count"),
		"How many messages did this router sent to this BGP peer (limited to IPv4).",
		[]string{"name"}, nil,
	)

	bgpPeerOutQueue = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "out_queue_count"),
		"PeerState.OutQ",
		[]string{"name"}, nil,
	)
	bgpPeerFlops = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "flop_count"),
		"PeerState.Flops",
		[]string{"name"}, nil,
	)
	bgpPeerSendCommunityFlag = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "send_community"),
		"PeerState.SendCommunity",
		[]string{"name"}, nil,
	)
	bgpPeerRemovePrivateAsFlag = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "remove_private_as"),
		"PeerState.RemovePrivateAs",
		[]string{"name"}, nil,
	)
	bgpPeerPasswodSetFlag = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "password_set"),
		"Whether the GoBGP peer has been configured (1) for authentication or not (0)",
		[]string{"name"}, nil,
	)
	bgpPeerType = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "peer", "type"),
		"PeerState.PeerType",
		[]string{"name"}, nil,
	)
)

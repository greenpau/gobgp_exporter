# GoBGP Exporter

Export GoBGP data to Prometheus.

![](./assets/docs/images/gobgp_exporter.png)


To run it:

```bash
cd $GOPATH/src
mkdir -p github.com/greenpau
cd github.com/greenpau
git clone https://github.com/greenpau/gobgp_exporter.git
cd gobgp_exporter
make
make qtest
```

## Exported Metrics

| **Metric** | **Description** | **Labels** |
| ------ | ------- | ------ |
| `gobgp_router_up` | Is GoBGP up and responds to queries (1) or is it down (0). | |
| `gobgp_router_id` | What is GoBGP router ID. | `id` |
| `gobgp_router_asn` | What is GoBGP AS number. | |
| `gobgp_router_failed_req_count` | The number of failed requests to GoBGP router. | |
| `gobgp_router_next_poll` | The timestamp of the next potential scrape of the router. | |
| `gobgp_router_scrape_time` | The amount of time it took to scrape the router. | |
| `gobgp_route_total_destination_count` | The number of routes on per address family and route table basis | `address_family`, `route_table`, `vrf_name` |
| `gobgp_route_total_path_count` | The number of available paths to destinations on per address family and route table basis | `address_family`, `route_table`, `vrf_name` |
| `gobgp_route_accepted_path_count` | The number of accepted paths to destinations on per address family and route table basis | `address_family`, `route_table`, `vrf_name` |
| `gobgp_peer_count` | The number of BGP peers | |
| `gobgp_peer_up` | Is the peer up and in established state (1) or it is not (0). | `name` |
| `gobgp_peer_asn` | What is the AS number of the peer | `name` |
| `gobgp_peer_local_asn` | What is the AS number presented to the peer by this router. | `name` |
| `gobgp_peer_admin_state` | Is the peer configured for being Up (0), Down (1), or PFX_CT (2) | `name` |
| `gobgp_peer_session_state` | What is the state of BGP session to the peer - unknown (0), idle (1), connect (2), active (3), opensent (4), openconfirm (5), established (6) | `name` |
| `gobgp_peer_received_message_total_count` | The total number of messages the BGP peer sent to this router (limited to IPv4). | `name` |
| `gobgp_peer_received_notification_message_count` | How many Notification messages did the BGP peer sent to this router (limited to IPv4). | `name` |
| `gobgp_peer_received_update_message_count` | How many Update messages did the BGP peer sent to this router (limited to IPv4). | `name` |
| `gobgp_peer_received_open_message_count` | How many Open messages did the BGP peer sent to this router (limited to IPv4). | `name` |
| `gobgp_peer_received_keepalive_message_count` | How many messages did the BGP peer sent to this router (limited to IPv4). | `name` |
| `gobgp_peer_received_refresh_message_count` | How many Refresh messages did the BGP peer sent to this router (limited to IPv4). | `name` |
| `gobgp_peer_received_withdraw_update_message_count` | How many WithdrawUpdate messages did the BGP peer sent to this router (limited to IPv4). | `name` |
| `gobgp_peer_received_withdraw_prefix_message_count` | How many messages did the BGP peer sent to this router (limited to IPv4). | `name` |
| `gobgp_peer_sent_message_total_count` | The total number of messages this router sent to this BGP peer (limited to IPv4). | `name` |
| `gobgp_peer_sent_notification_message_count` | How many Notification messages did this router sent to this BGP peer (limited to IPv4). | `name` |
| `gobgp_peer_sent_update_message_count` | How many Update messages did this router sent to this BGP peer (limited to IPv4). | `name` |
| `gobgp_peer_sent_open_message_count` | How many Open messages did this router sent to this BGP peer (limited to IPv4). | `name` |
| `gobgp_peer_sent_keepalive_message_count` | How many messages did this router sent to this BGP peer (limited to IPv4). | `name` |
| `gobgp_peer_sent_refresh_message_count` | How many Refresh messages did this router sent to this BGP peer (limited to IPv4). | `name` |
| `gobgp_peer_sent_withdraw_update_message_count` | How many WithdrawUpdate messages did this router sent to this BGP peer (limited to IPv4). | `name` |
| `gobgp_peer_sent_withdraw_prefix_message_count` | How many messages did this router sent to this BGP peer (limited to IPv4). | `name` |
| `gobgp_peer_out_queue_count` | PeerState.OutQ | `name` |
| `gobgp_peer_flop_count` | PeerState.Flops | `name` |
| `gobgp_peer_send_community` | PeerState.SendCommunity | `name` |
| `gobgp_peer_remove_private_as` | PeerState.RemovePrivateAs | `name` |
| `gobgp_peer_password_set` | Whether the GoBGP peer has been configured (1) for authentication or not (0) | `name` |
| `gobgp_peer_type` | PeerState.PeerType | `name` |

For example:

```
# HELP gobgp_peer_admin_state Is the peer configured for being Up (0), Down (1), or PFX_CT (2)
# TYPE gobgp_peer_admin_state gauge
gobgp_peer_admin_state{name="10.0.2.100"} 0
# HELP gobgp_peer_asn What is the AS number of the peer
# TYPE gobgp_peer_asn gauge
gobgp_peer_asn{name="10.0.2.100"} 65001
# HELP gobgp_peer_count The number of BGP peers
# TYPE gobgp_peer_count gauge
gobgp_peer_count 1
# HELP gobgp_peer_flop_count PeerState.Flops
# TYPE gobgp_peer_flop_count gauge
gobgp_peer_flop_count{name="10.0.2.100"} 0
# HELP gobgp_peer_local_asn What is the AS number presented to the peer by this router.
# TYPE gobgp_peer_local_asn gauge
gobgp_peer_local_asn{name="10.0.2.100"} 0
# HELP gobgp_peer_out_queue_count PeerState.OutQ
# TYPE gobgp_peer_out_queue_count gauge
gobgp_peer_out_queue_count{name="10.0.2.100"} 0
# HELP gobgp_peer_password_set Whether the GoBGP peer has been configured (1) for authentication or not (0)
# TYPE gobgp_peer_password_set gauge
gobgp_peer_password_set{name="10.0.2.100"} 0
# HELP gobgp_peer_received_keepalive_message_count How many messages did the BGP peer sent to this router (limited to IPv4).
# TYPE gobgp_peer_received_keepalive_message_count gauge
gobgp_peer_received_keepalive_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_received_message_total_count The total number of messages the BGP peer sent to this router (limited to IPv4).
# TYPE gobgp_peer_received_message_total_count gauge
gobgp_peer_received_message_total_count{name="10.0.2.100"} 0
# HELP gobgp_peer_received_notification_message_count How many Notification messages did the BGP peer sent to this router (limited to IPv4).
# TYPE gobgp_peer_received_notification_message_count gauge
gobgp_peer_received_notification_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_received_open_message_count How many Open messages did the BGP peer sent to this router (limited to IPv4).
# TYPE gobgp_peer_received_open_message_count gauge
gobgp_peer_received_open_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_received_refresh_message_count How many Refresh messages did the BGP peer sent to this router (limited to IPv4).
# TYPE gobgp_peer_received_refresh_message_count gauge
gobgp_peer_received_refresh_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_received_update_message_count How many Update messages did the BGP peer sent to this router (limited to IPv4).
# TYPE gobgp_peer_received_update_message_count gauge
gobgp_peer_received_update_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_received_withdraw_prefix_message_count How many messages did the BGP peer sent to this router (limited to IPv4).
# TYPE gobgp_peer_received_withdraw_prefix_message_count gauge
gobgp_peer_received_withdraw_prefix_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_received_withdraw_update_message_count How many WithdrawUpdate messages did the BGP peer sent to this router (limited to IPv4).
# TYPE gobgp_peer_received_withdraw_update_message_count gauge
gobgp_peer_received_withdraw_update_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_remove_private_as PeerState.RemovePrivateAs
# TYPE gobgp_peer_remove_private_as gauge
gobgp_peer_remove_private_as{name="10.0.2.100"} 0
# HELP gobgp_peer_send_community PeerState.SendCommunity
# TYPE gobgp_peer_send_community gauge
gobgp_peer_send_community{name="10.0.2.100"} 0
# HELP gobgp_peer_sent_keepalive_message_count How many messages did this router sent to this BGP peer (limited to IPv4).
# TYPE gobgp_peer_sent_keepalive_message_count gauge
gobgp_peer_sent_keepalive_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_sent_message_total_count The total number of messages this router sent to this BGP peer (limited to IPv4).
# TYPE gobgp_peer_sent_message_total_count gauge
gobgp_peer_sent_message_total_count{name="10.0.2.100"} 0
# HELP gobgp_peer_sent_notification_message_count How many Notification messages did this router sent to this BGP peer (limited to IPv4).
# TYPE gobgp_peer_sent_notification_message_count gauge
gobgp_peer_sent_notification_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_sent_open_message_count How many Open messages did this router sent to this BGP peer (limited to IPv4).
# TYPE gobgp_peer_sent_open_message_count gauge
gobgp_peer_sent_open_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_sent_refresh_message_count How many Refresh messages did this router sent to this BGP peer (limited to IPv4).
# TYPE gobgp_peer_sent_refresh_message_count gauge
gobgp_peer_sent_refresh_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_sent_update_message_count How many Update messages did this router sent to this BGP peer (limited to IPv4).
# TYPE gobgp_peer_sent_update_message_count gauge
gobgp_peer_sent_update_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_sent_withdraw_prefix_message_count How many messages did this router sent to this BGP peer (limited to IPv4).
# TYPE gobgp_peer_sent_withdraw_prefix_message_count gauge
gobgp_peer_sent_withdraw_prefix_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_sent_withdraw_update_message_count How many WithdrawUpdate messages did this router sent to this BGP peer (limited to IPv4).
# TYPE gobgp_peer_sent_withdraw_update_message_count gauge
gobgp_peer_sent_withdraw_update_message_count{name="10.0.2.100"} 0
# HELP gobgp_peer_session_state What is the state of BGP session to the peer - unknown (0), idle (1), connect (2), active (3), opensent (4), openconfirm (5), established (6)
# TYPE gobgp_peer_session_state gauge
gobgp_peer_session_state{name="10.0.2.100"} 3
# HELP gobgp_peer_type PeerState.PeerType
# TYPE gobgp_peer_type gauge
gobgp_peer_type{name="10.0.2.100"} 0
# HELP gobgp_peer_up Is the peer up and in established state (1) or it is not (0).
# TYPE gobgp_peer_up gauge
gobgp_peer_up{name="10.0.2.100"} 0
# HELP gobgp_route_accepted_path_count The number of accepted paths to destinations on per address family and route table basis
# TYPE gobgp_route_accepted_path_count gauge
gobgp_route_accepted_path_count{address_family="evpn",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="evpn",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4_encap",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4_encap",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4_mpls",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4_mpls",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4_vpn",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4_vpn",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4_vpn_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv4_vpn_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6_encap",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6_encap",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6_mpls",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6_mpls",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6_vpn",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6_vpn",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6_vpn_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="ipv6_vpn_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="l2_vpn_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_accepted_path_count{address_family="l2_vpn_flowspec",route_table="local",vrf_name="default"} 0
# HELP gobgp_route_total_destination_count The number of routes on per address family and route table basis
# TYPE gobgp_route_total_destination_count gauge
gobgp_route_total_destination_count{address_family="evpn",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="evpn",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv4",route_table="global",vrf_name="default"} 3
gobgp_route_total_destination_count{address_family="ipv4",route_table="local",vrf_name="default"} 3
gobgp_route_total_destination_count{address_family="ipv4_encap",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv4_encap",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv4_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv4_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv4_mpls",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv4_mpls",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv4_vpn",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv4_vpn",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv4_vpn_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv4_vpn_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6_encap",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6_encap",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6_mpls",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6_mpls",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6_vpn",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6_vpn",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6_vpn_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="ipv6_vpn_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="l2_vpn_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_total_destination_count{address_family="l2_vpn_flowspec",route_table="local",vrf_name="default"} 0
# HELP gobgp_route_total_path_count The number of available paths to destinations on per address family and route table basis
# TYPE gobgp_route_total_path_count gauge
gobgp_route_total_path_count{address_family="evpn",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="evpn",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv4",route_table="global",vrf_name="default"} 3
gobgp_route_total_path_count{address_family="ipv4",route_table="local",vrf_name="default"} 3
gobgp_route_total_path_count{address_family="ipv4_encap",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv4_encap",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv4_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv4_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv4_mpls",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv4_mpls",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv4_vpn",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv4_vpn",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv4_vpn_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv4_vpn_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6_encap",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6_encap",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6_mpls",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6_mpls",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6_vpn",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6_vpn",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6_vpn_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="ipv6_vpn_flowspec",route_table="local",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="l2_vpn_flowspec",route_table="global",vrf_name="default"} 0
gobgp_route_total_path_count{address_family="l2_vpn_flowspec",route_table="local",vrf_name="default"} 0
# HELP gobgp_router_asn What is GoBGP AS number.
# TYPE gobgp_router_asn gauge
gobgp_router_asn 65001
# HELP gobgp_router_failed_req_count The number of failed requests to GoBGP router.
# TYPE gobgp_router_failed_req_count counter
gobgp_router_failed_req_count 0
# HELP gobgp_router_id What is GoBGP router ID.
# TYPE gobgp_router_id gauge
gobgp_router_id 1
# HELP gobgp_router_next_poll The timestamp of the next potential scrape of the router.
# TYPE gobgp_router_next_poll counter
gobgp_router_next_poll 0
# HELP gobgp_router_scrape_time The amount of time it took to scrape the router.
# TYPE gobgp_router_scrape_time gauge
gobgp_router_scrape_time 0.019395555
# HELP gobgp_router_up Is GoBGP up and responds to queries (1) or is it down (0).
# TYPE gobgp_router_up gauge
gobgp_router_up 1
```

## Flags

```bash
$ bin/gobgp-exporter --help

gobgp-exporter - Prometheus Exporter for GoBGP

Usage: gobgp-exporter [arguments]

  -auth.token string
        The X-Token for accessing the exporter itself (default "anonymous")
  -gobgp.address string
        gRPC API address of GoBGP server. (default "127.0.0.1:50051")
  -gobgp.poll-interval int
        The minimum interval (in seconds) between collections from a GoBGP server. (default 15)
  -gobgp.timeout int
        Timeout on gRPC requests to a GoBGP server. (default 2)
  -gobgp.tls
        Whether to enable TLS for gRPC API access.
  -gobgp.tls-ca string
        Optional path to PEM file with CA certificates to be trusted for gRPC API access.
  -gobgp.tls-client-cert string
        Optional path to PEM file with client certificate to be used for client authentication.
  -gobgp.tls-client-key string
        Optional path to PEM file with client key to be used for client authentication.
  -gobgp.tls-server-name string
        Optional hostname to verify API server as.
  -log.level string
        logging severity level (default "info")
  -metrics
        Display available metrics
  -version
        version information
  -web.listen-address string
        Address to listen on for web interface and telemetry. (default ":9474")
  -web.telemetry-path string
        Path under which to expose metrics. (default "/metrics")

Documentation: https://github.com/greenpau/gobgp_exporter/
```

* __`gobgp.address`:__ Address (host and port) of the GoBGP instance we should
    connect to. This could be a local GoBGP server (`127.0.0.0:50051`, for
    instance), or the address of a remote GoBGP server.
* __`gobgp.tls`:__ Enable TLS for the GoBGP connection. (default: false)
* __`gobgp.tls-ca`:__ Optional path to a PEM file containing certificate authorities to verify GoBGP server certificate against. If empty, the host's root CA set is used instead. (default: empty)
* __`gobgp.tls-client-cert`:__ Optional path to a PEM file containing the client certificate to authenticate with. (default: empty)
* __`gobgp.tls-client-key`:__ Optional path to a PEM file containing the key for theclient certificate to authenticate with. (default: empty)
* __`gobgp.tls-server-name`:__ Optional server name to verify GoBGP server certificate against. If empty, verification will be using the hostname or IP used in `gobgp.address`. (default: empty)
* __`gobgp.timeout`:__ Timeout on gRPC requests to GoBGP.
* __`gobgp.poll-interval`:__ The minimum interval (in seconds) between collections from GoBGP server. (default: 15 seconds)
* __`gobgp.peers`:__ The file containing the mapping between `router_id` and the name (e.g. `hostname`) of a remote peer.
* __`auth.token`:__ Enable X-Token authentication for accessing the exporter itself.
* __`version`:__ Show application version.
* __`web.listen-address`:__ Address to listen on for web interface and telemetry.
* __`web.telemetry-path`:__ Path under which to expose metrics.

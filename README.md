# GoBGP Exporter

Export GoBGP data to Prometheus.

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

| Metric | Meaning | Labels |
| ------ | ------- | ------ |
| `gobgp_up` | Is GoBGP up and responds to queries (1) or is it down (0). | |
| `gobgp_asn` | What is GoBGP router ID and AS number. | `router_id` |
| `gobgp_connected_at` | When was the last successful connection to GoBGP. | `router_id` |
| `gobgp_lost_connection_at` | When did the exporter lose connection to GoBGP router. | `router_id` |
| `gobgp_failed_query_count` | The number of failed queries to GoBGP router. | `router_id` |
| `gobgp_next_poll` | The timestamp of the next potential poll of GoBGP server. | `router_id` |
| `gobgp_route_count` | The number of routes on per address family and resource type basis. | `router_id`, `address_family`, `resource_type` |
| `gobgp_peer_count` | The number of BGP peers | `router_id` |
| `gobgp_peer_name` | **TODO** The name associated with a remote peer. The names must provided to the exporter in the form of YAML file. | `router_id` |
| `gobgp_peer_up` | Is GoBGP peer up or down (0). | `router_id`, `peer_router_id` |
| `gobgp_peer_asn` | What is AS number for a peer. | `router_id`, `peer_router_id` |
| `gobgp_peer_admin_state` | Is the peer configured for being Up (0), Down (1), or PFX_CT (2) | `router_id`, `peer_router_id` |
| `gobgp_peer_session_state` | What is the state of BGP session to the peer: unknown (0), idle (1), connect (2), active (3), opensent (4), openconfirm (5), established (6) | `router_id`, `peer_router_id` |
| `gobgp_peer_received_route_count` | How many routes did the BGP peer sent to this router (limited to IPv4). | `router_id`, `peer_router_id` |
| `gobgp_peer_accepted_route_count` | How many routes were accepted from the routes received from this BGP peer (limited to IPv4). | `router_id`, `peer_router_id` |
| `gobgp_peer_advertised_route_count` | How many routes were advertised to this BGP peer (limited to IPv4). | `router_id`, `peer_router_id` |
| `gobgp_peer_out_queue_count` | `PeerState.OutQ` | `router_id`, `peer_router_id` |
| `gobgp_peer_flop_count` | `PeerState.Flops` | `router_id`, `peer_router_id` |
| `gobgp_peer_send_community` | `PeerState.SendCommunity` | `router_id`, `peer_router_id` |
| `gobgp_peer_remove_private_as` | `PeerState.RemovePrivateAs`: None (0), All (1), Replace (2) | `router_id`, `peer_router_id` |
| `gobgp_peer_password_set` | **TODO** `PeerState.AuthPassword` | `router_id`, `peer_router_id` |
| `gobgp_peer_type` | `PeerState.PeerType`: internal (0), external (1) | `router_id`, `peer_router_id` |

For example:

```
gobgp_asn{router_id="192.168.56.3"} 65001
gobgp_connected_at{router_id="192.168.56.3"} 1.534704149e+09
gobgp_exporter_build_info{branch="master",goversion="go1.10.2",revision="687ae72-dirty",version="1.0.0"} 1
gobgp_lost_connection_at{router_id="192.168.56.3"} 0
gobgp_failed_query_count{router_id="192.168.56.3"} 2
gobgp_route_count{address_family="evpn",resource_type="global",router_id="192.168.56.3"} 4
gobgp_route_count{address_family="evpn",resource_type="local",router_id="192.168.56.3"} 4
gobgp_route_count{address_family="ipv4",resource_type="global",router_id="192.168.56.3"} 1
gobgp_route_count{address_family="ipv4",resource_type="local",router_id="192.168.56.3"} 1
gobgp_peer_count{router_id="192.168.56.2"} 2
gobgp_peer_up{peer_router_id="192.168.56.4",router_id="192.168.56.3"} 1
gobgp_peer_up{peer_router_id="192.168.56.5",router_id="192.168.56.3"} 1
gobgp_peer_asn{peer_router_id="192.168.56.4",router_id="192.168.56.3"} 65001
gobgp_peer_asn{peer_router_id="192.168.56.5",router_id="192.168.56.3"} 65001
gobgp_peer_admin_state{peer_router_id="192.168.56.4",router_id="192.168.56.3"} 0
gobgp_peer_admin_state{peer_router_id="192.168.56.5",router_id="192.168.56.3"} 0
gobgp_up 1
```

## Flags

```bash
./gobgp_exporter --help
```

* __`gobgp.address`:__ Address (host and port) of the GoBGP instance we should
    connect to. This could be a local GoBGP server (`127.0.0.0:50051`, for
    instance), or the address of a remote GoBGP server.
* __`gobgp.timeout`:__ Timeout on gRPC requests to GoBGP.
* __`gobgp.poll-interval`:__ The minimum interval (in seconds) between collections from GoBGP server. (default: 15 seconds)
* __`gobgp.peers`:__ The file containing the mapping between `router_id` and the name (e.g. `hostname`) of a remote peer.
* __`version`:__ Show application version.
* __`web.listen-address`:__ Address to listen on for web interface and telemetry.
* __`web.telemetry-path`:__ Path under which to expose metrics.

## Outstanding Issues

The following is a list of issues related to GoBGP package itself:
- `gobgp_peer_session_state` reports `0`, although a router maybe in `established` (`5`) state.
  The metric is being derived from `PeerState.SessionState`

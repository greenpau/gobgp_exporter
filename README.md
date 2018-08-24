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

For example:

```
gobgp_asn{router_id="192.168.56.3"} 65001
gobgp_connected_at{router_id="192.168.56.3"} 1.534704149e+09
gobgp_exporter_build_info{branch="master",goversion="go1.10.2",revision="687ae72-dirty",version="1.0.0"} 1
gobgp_lost_connection_at{router_id="192.168.56.3"} 0
gobgp_failed_query_count{router_id="192.168.56.3"} 2
gobgp_route_count{address_family="evpn",resource_type="global",router_id="192.168.56.2"} 4
gobgp_route_count{address_family="evpn",resource_type="local",router_id="192.168.56.2"} 4
gobgp_route_count{address_family="ipv4",resource_type="global",router_id="192.168.56.2"} 1
gobgp_route_count{address_family="ipv4",resource_type="local",router_id="192.168.56.2"} 1
gobgp_up 1
```

### Flags

```bash
./gobgp_exporter --help
```

* __`gobgp.address`:__ Address (host and port) of the GoBGP instance we should
    connect to. This could be a local GoBGP server (`127.0.0.0:50051`, for
    instance), or the address of a remote GoBGP server.
* __`gobgp.timeout`:__ Timeout on gRPC requests to GoBGP.
* __`gobgp.poll-interval`:__ The minimum interval (in seconds) between collections from GoBGP server. (default: 15 seconds)
* __`version`:__ Show application version.
* __`web.listen-address`:__ Address to listen on for web interface and telemetry.
* __`web.telemetry-path`:__ Path under which to expose metrics.

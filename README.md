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

For example:

```
gobgp_asn{router_id="192.168.56.3"} 65001
gobgp_connected_at{router_id="192.168.56.3"} 1.534704149e+09
gobgp_exporter_build_info{branch="master",goversion="go1.10.2",revision="687ae72-dirty",version="1.0.0"} 1
gobgp_lost_connection_at{router_id="192.168.56.3"} 0
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
* __`version`:__ Show application version.
* __`web.listen-address`:__ Address to listen on for web interface and telemetry.
* __`web.telemetry-path`:__ Path under which to expose metrics.

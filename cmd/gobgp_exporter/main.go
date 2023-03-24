package main

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-kit/log/level"
	exporter "github.com/greenpau/gobgp_exporter/pkg/gobgp_exporter"
	"github.com/prometheus/common/promlog"
)

func loadCertificatePEM(filePath string) (*x509.Certificate, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	rest := content
	var block *pem.Block
	var cert *x509.Certificate
	for len(rest) > 0 {
		block, rest = pem.Decode(content)
		if block == nil {
			// no PEM data found, rest will not have been modified
			break
		}
		content = rest
		switch block.Type {
		case "CERTIFICATE":
			cert, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			return cert, err
		default:
			// not the PEM block we're looking for
			continue
		}
	}
	return nil, errors.New("no certificate PEM block found")
}

func loadKeyPEM(filePath string) (crypto.PrivateKey, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	rest := content
	var block *pem.Block
	var key crypto.PrivateKey
	for len(rest) > 0 {
		block, rest = pem.Decode(content)
		if block == nil {
			// no PEM data found, rest will not have been modified
			break
		}
		switch block.Type {
		case "RSA PRIVATE KEY":
			key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			return key, err
		case "PRIVATE KEY":
			key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			return key, err
		case "EC PRIVATE KEY":
			key, err = x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			return key, err
		default:
			// not the PEM block we're looking for
			continue
		}
	}
	return nil, errors.New("no private key PEM block found")
}

func main() {
	var listenAddress string
	var metricsPath string
	var serverAddress string
	var serverTLS bool
	var serverTLSCAPath string
	var serverTLSServerName string
	var serverTLSClientCertPath string
	var serverTLSClientKeyPath string
	var pollTimeout int
	var pollInterval int
	var isShowMetrics bool
	var isShowVersion bool
	var logLevel string
	var authToken string

	flag.StringVar(&listenAddress, "web.listen-address", ":9474", "Address to listen on for web interface and telemetry.")
	flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	flag.StringVar(&serverAddress, "gobgp.address", "127.0.0.1:50051", "gRPC API address of GoBGP server.")
	flag.BoolVar(&serverTLS, "gobgp.tls", false, "Whether to enable TLS for gRPC API access.")
	flag.StringVar(&serverTLSCAPath, "gobgp.tls-ca", "", "Optional path to PEM file with CA certificates to be trusted for gRPC API access.")
	flag.StringVar(&serverTLSServerName, "gobgp.tls-server-name", "", "Optional hostname to verify API server as.")
	flag.StringVar(&serverTLSClientCertPath, "gobgp.tls-client-cert", "", "Optional path to PEM file with client certificate to be used for client authentication.")
	flag.StringVar(&serverTLSClientKeyPath, "gobgp.tls-client-key", "", "Optional path to PEM file with client key to be used for client authentication.")
	flag.IntVar(&pollTimeout, "gobgp.timeout", 2, "Timeout on gRPC requests to a GoBGP server.")
	flag.IntVar(&pollInterval, "gobgp.poll-interval", 15, "The minimum interval (in seconds) between collections from a GoBGP server.")
	flag.StringVar(&authToken, "auth.token", "anonymous", "The X-Token for accessing the exporter itself")
	flag.BoolVar(&isShowMetrics, "metrics", false, "Display available metrics")
	flag.BoolVar(&isShowVersion, "version", false, "version information")
	flag.StringVar(&logLevel, "log.level", "info", "logging severity level")

	usageHelp := func() {
		fmt.Fprintf(os.Stderr, "\n%s - Prometheus Exporter for GoBGP\n\n", exporter.GetExporterName())
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments]\n\n", exporter.GetExporterName())
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDocumentation: https://github.com/greenpau/gobgp_exporter/\n\n")
	}
	flag.Usage = usageHelp
	flag.Parse()

	opts := exporter.Options{
		Address: serverAddress,
		Timeout: pollTimeout,
	}

	allowedLogLevel := &promlog.AllowedLevel{}
	if err := allowedLogLevel.Set(logLevel); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	promlogConfig := &promlog.Config{
		Level: allowedLogLevel,
	}

	logger := promlog.New(promlogConfig)
	opts.Logger = logger

	if serverTLS {
		opts.TLS = new(tls.Config)
		if len(serverTLSCAPath) > 0 {
			// assuming PEM file here
			pemCerts, err := os.ReadFile(filepath.Clean(serverTLSCAPath))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not read TLS CA PEM file %q: %s\n", serverTLSCAPath, err)
				os.Exit(1)
			}

			opts.TLS.RootCAs = x509.NewCertPool()
			ok := opts.TLS.RootCAs.AppendCertsFromPEM(pemCerts)
			if !ok {
				fmt.Fprintf(os.Stderr, "Could not parse any TLS CA certificate from PEM file %q: %s\n", serverTLSCAPath, err)
				os.Exit(1)
			}
		}
		if len(serverTLSServerName) > 0 {
			opts.TLS.ServerName = serverTLSServerName
		}
		if len(serverTLSClientCertPath) > 0 && len(serverTLSClientKeyPath) > 0 {
			// again assuming PEM file
			cert, err := loadCertificatePEM(serverTLSClientCertPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load client certificate: %s\n", err)
			}
			key, err := loadKeyPEM(serverTLSClientKeyPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load client key: %s\n", err)
			}
			opts.TLS.Certificates = []tls.Certificate{
				{
					Certificate: [][]byte{cert.Raw},
					PrivateKey:  key,
				},
			}
		} else if len(serverTLSClientCertPath) > 0 || len(serverTLSClientKeyPath) > 0 {
			fmt.Fprintln(os.Stderr, "Only one of client certificate and key was set, must set both.")
			os.Exit(1)
		}
	}

	if isShowVersion {
		fmt.Fprintf(os.Stdout, "%s %s", exporter.GetExporterName(), exporter.GetVersion())
		if exporter.GetRevision() != "" {
			fmt.Fprintf(os.Stdout, ", commit: %s\n", exporter.GetRevision())
		} else {
			fmt.Fprint(os.Stdout, "\n")
		}
		os.Exit(0)
	}

	if isShowMetrics {
		e := &exporter.RouterNode{}
		fmt.Fprintf(os.Stdout, "%s\n", e.GetMetricsTable())
		os.Exit(0)
	}

	level.Info(logger).Log(
		"msg", "Starting exporter",
		"exporter", exporter.GetExporterName(),
		"version", exporter.GetVersionInfo(),
		"build_context", exporter.GetVersionBuildContext(),
	)

	e, err := exporter.NewExporter(opts)
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to init properly",
			"error", err.Error(),
		)
		os.Exit(1)
	}

	e.SetPollInterval(int64(pollInterval))
	if err := e.AddAuthenticationToken(authToken); err != nil {
		level.Error(logger).Log(
			"msg", "failed to add authentication token",
			"error", err.Error(),
		)
		os.Exit(1)
	}

	level.Info(logger).Log(
		"msg", "exporter configuration",
		"min_scrape_interval", e.GetPollInterval(),
	)

	http.HandleFunc(metricsPath, func(w http.ResponseWriter, r *http.Request) {
		e.Scrape(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		e.Summary(metricsPath, w, r)
	})

	level.Info(logger).Log("listen_on ", listenAddress)

	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		level.Error(logger).Log(
			"msg", "listener failed",
			"error", err.Error(),
		)
		os.Exit(1)
	}
}

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
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/common/version"
)

const (
	namespace = "gobgp"
)

var (
	appName    = "gobgp-exporter"
	appVersion = "[untracked]"
	gitBranch  string
	gitCommit  string
	buildUser  string // whoami
	buildDate  string // date -u
)

// Exporter collects GoBGP data from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
	sync.RWMutex
	timeout      int
	address      string
	pollInterval int64
	Node         *RouterNode
	Tokens       map[string]bool
	logger       log.Logger
}

// Options are the options for the initialization of an instance of the
// Exporter.
type Options struct {
	Address string
	TLS     *tls.Config
	Timeout int
	Logger  log.Logger
}

// NewExporter returns an initialized Exporter.
func NewExporter(opts Options) (*Exporter, error) {
	version.Version = appVersion
	version.Revision = gitCommit
	version.Branch = gitBranch
	version.BuildUser = buildUser
	version.BuildDate = buildDate
	e := Exporter{
		timeout: opts.Timeout,
		address: opts.Address,
		Tokens:  make(map[string]bool),
		logger:  opts.Logger,
	}

	n, err := NewRouterNode(opts.Address, opts.Timeout, opts.TLS, opts.Logger)
	if err != nil {
		return nil, err
	}
	e.Node = n
	level.Debug(e.logger).Log(
		"msg", "NewExporter() initialized successfully",
	)

	return &e, nil
}

// GetVersionInfo returns exporter info.
func GetVersionInfo() string {
	return version.Info()
}

// GetVersionBuildContext returns exporter build context.
func GetVersionBuildContext() string {
	return version.BuildContext()
}

// GetVersion returns exporter version.
func GetVersion() string {
	return version.Version
}

// GetRevision returns exporter revision.
func GetRevision() string {
	return version.Revision
}

// GetExporterName returns exporter name.
func GetExporterName() string {
	return appName
}

// SetPollInterval sets exporter's minimal polling/scraping interval.
func (e *Exporter) SetPollInterval(i int64) {
	e.pollInterval = i
	if e.Node.pollInterval == 0 {
		e.Node.pollInterval = i
	}
}

// GetPollInterval returns exporters minimal polling/scraping interval.
func (e *Exporter) GetPollInterval() int64 {
	return e.pollInterval
}

// Scrape scrapes individual nodes.
func (e *Exporter) Scrape(w http.ResponseWriter, r *http.Request) {
	if _, authorized := e.authorize(r); !authorized {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	level.Debug(e.logger).Log(
		"msg", "calls Scrape()",
	)

	start := time.Now()
	registry := prometheus.NewRegistry()
	registry.MustRegister(e.Node)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	duration := time.Since(start).Seconds()
	level.Debug(e.logger).Log(
		"msg", "completed Scrape()",
		"took", duration,
	)
}

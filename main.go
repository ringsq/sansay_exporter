// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const targetPath = "/SSConfig/webresources/stats/"

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9116").String()
	dryRun        = kingpin.Flag("dry-run", "Only verify configuration is valid and exit.").Default("false").Bool()

	// Metrics about the sansay exporter itself.
	sansayDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "sansay_collection_duration_seconds",
			Help: "Duration of collections by the sansay exporter",
		},
	)
	sansayRequestErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sansay_request_errors_total",
			Help: "Errors in requests to the sansay exporter",
		},
	)
)

func init() {
	prometheus.MustRegister(sansayDuration)
	prometheus.MustRegister(sansayRequestErrors)
	prometheus.MustRegister(version.NewCollector("sansay_exporter"))
}

func handler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	useSoap := false
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "'target' parameter must be specified", 400)
		sansayRequestErrors.Inc()
		return
	}
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	protocol := r.URL.Query().Get("protocol")
	if protocol == "" {
		protocol = "https"
	}
	api := r.URL.Query().Get("api")
	if strings.ToLower(api) == "soap" {
		useSoap = true
	}

	logger = log.With(logger, "target", target)
	level.Debug(logger).Log("msg", "Starting scrape", "module")

	start := time.Now()
	registry := prometheus.NewRegistry()
	collector := collector{target: fmt.Sprintf("%s://%s", protocol, target), targetPath: targetPath, useSoap: useSoap, username: username, password: password, logger: logger}
	registry.MustRegister(collector)
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	duration := time.Since(start).Seconds()
	sansayDuration.Observe(duration)
	level.Debug(logger).Log("msg", "Finished scrape", "duration_seconds", duration)
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting sansay_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", version.BuildContext())

	// Exit if in dry-run mode.
	if *dryRun {
		level.Info(logger).Log("msg", "Configuration parsed successfully")
		return
	}

	http.Handle("/metrics", promhttp.Handler()) // Normal metrics endpoint for sansay exporter itself.
	// Endpoint to do sansay scrapes.
	http.HandleFunc("/sansay", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, logger)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head>
            <title>Sansay Exporter</title>
            <style>
            label{
            display:inline-block;
            width:75px;
            }
            form label {
            margin: 10px;
            }
            form input {
            margin: 10px;
            }
            </style>
            </head>
            <body>
            <h1>Sansay Exporter</h1>
            <form action="/sansay">
            <label>Target:</label> <input type="text" name="target" placeholder="X.X.X.X" value="1.2.3.4"><br>
            <input type="submit" value="Submit">
            </form>
            </body>
            </html>`))
	})

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}

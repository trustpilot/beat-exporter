package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"beat-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

type Params struct {
	ServiceName   string
	ListenAddress string
	TLSCertFile   string
	TLSKeyFile    string
	MetricsPath   string
	BeatURI       string
	BeatTimeout   time.Duration
	SystemBeat    bool
}

func NewParams(
	listenAddress,
	tlsCertFile,
	tlsKeyFile,
	metricsPath,
	beatURI string,
	BeatTimeout time.Duration,
	systemBeat bool) *Params {
	return &Params{
		ServiceName:   "beat_exporter",
		ListenAddress: listenAddress,
		TLSCertFile:   tlsCertFile,
		TLSKeyFile:    tlsKeyFile,
		MetricsPath:   metricsPath,
		BeatURI:       beatURI,
		BeatTimeout:   BeatTimeout,
		SystemBeat:    systemBeat,
	}
}

type Exporter struct {
	ServiceParams      *Params
	BeatURL            *url.URL
	HTTPClient         *http.Client
	BeatInfo           *collector.BeatInfo
	PrometheusRegistry *prometheus.Registry
}

func NewExporter(params *Params) *Exporter {
	exporter := &Exporter{ServiceParams: params}
	beatURL, err := url.Parse(exporter.ServiceParams.BeatURI)
	if err != nil {
		log.Fatalf("failed to parse beat.uri, error: %v", err)
		panic(err)
	}
	exporter.BeatURL = beatURL
	exporter.HTTPClient = &http.Client{
		Timeout: exporter.ServiceParams.BeatTimeout,
	}

	exporter.handleURL()
	exporter.loadBeatInfo()
	exporter.loadPrometheusRegistry()
	return exporter
}

func (e *Exporter) handleURL() {
	if e.BeatURL.Scheme == "unix" {
		unixPath := e.BeatURL.Path
		e.BeatURL.Scheme = "http"
		e.BeatURL.Host = "localhost"
		e.BeatURL.Path = ""
		e.HTTPClient.Transport = &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", unixPath)
			},
		}
	}
}

func (e *Exporter) loadBeatInfo() {
	response, err := e.HTTPClient.Get(e.BeatURL.String())
	if err != nil {
		log.Errorf("Could not load beat type, with error: %v", err)
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Errorf("Beat URL: %q status code: %d", e.BeatURL.String(), response.StatusCode)
		panic(err)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Can't read body of response")
		panic(err)
	}

	beatInfo := &collector.BeatInfo{}
	if err := json.Unmarshal(bodyBytes, &beatInfo); err != nil {
		log.Error("Could not parse JSON response for target")
		panic(err)
	}

	log.WithFields(
		log.Fields{
			"beat":     beatInfo.Beat,
			"version":  beatInfo.Version,
			"name":     beatInfo.Name,
			"hostname": beatInfo.Hostname,
			"uuid":     beatInfo.UUID,
		}).Info("Target beat configuration loaded successfully!")

	e.BeatInfo = beatInfo
}

func (e *Exporter) loadPrometheusRegistry() {
	params := e.ServiceParams
	// version metric
	e.PrometheusRegistry = prometheus.NewRegistry()
	versionMetric := version.NewCollector(params.ServiceName)
	mainCollector := collector.NewMainCollector(e.HTTPClient, e.BeatURL, params.ServiceName, e.BeatInfo, params.SystemBeat)
	e.PrometheusRegistry.MustRegister(versionMetric)
	e.PrometheusRegistry.MustRegister(mainCollector)
}

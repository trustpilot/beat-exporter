package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/trustpilot/beat-exporter/collector"
)

func main() {
	var (
		Name          = "beat_exporter"
		listenAddress = flag.String("web.listen-address", ":9479", "Address to listen on for web interface and telemetry.")
		tlsCertFile   = flag.String("tls.certfile", "", "TLS certs file if you want to use tls instead of http")
		tlsKeyFile    = flag.String("tls.keyfile", "", "TLS key file if you want to use tls instead of http")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		beatURI       = flag.String("beat.uri", "http://localhost:5066", "HTTP API address of beat.")
		beatTimeout   = flag.Duration("beat.timeout", 10*time.Second, "Timeout for trying to get stats from beat.")
		showVersion   = flag.Bool("version", false, "Show version and exit")
	)
	flag.Parse()

	if *showVersion {
		fmt.Print(version.Print(Name))
		os.Exit(0)
	}

	log.SetLevel(log.DebugLevel)

	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyMsg: "message",
		},
	})

	beatURL, err := url.Parse(*beatURI)

	if err != nil {
		log.Fatalf("failed to parse beat.uri, error: %v", err)
	}

	httpClient := &http.Client{
		Timeout: *beatTimeout,
	}

	log.Info("Exploring target for beat type")

	var beatInfo *collector.BeatInfo

	for {
		beatInfo, err = loadBeatType(httpClient, *beatURL)
		if err != nil {
			log.Errorf("Could not load beat type, with error: %v, retrying in 5s", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	// version metric
	registry := prometheus.NewRegistry()
	versionMetric := version.NewCollector(Name)
	mainCollector := collector.NewMainCollector(httpClient, beatURL, Name, beatInfo)
	registry.MustRegister(versionMetric)
	registry.MustRegister(mainCollector)

	http.Handle(*metricsPath, promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			ErrorLog:           log.New(),
			DisableCompression: false,
			ErrorHandling:      promhttp.ContinueOnError}),
	)

	http.HandleFunc("/", IndexHandler(*metricsPath))

	log.WithFields(log.Fields{
		"addr": *listenAddress,
	}).Infof("Starting exporter with configured type: %s", beatInfo.Beat)

	if *tlsCertFile != "" && *tlsKeyFile != "" {
		if err := http.ListenAndServeTLS(*listenAddress, *tlsCertFile, *tlsKeyFile, nil); err != nil {

			log.WithFields(log.Fields{
				"err": err,
			}).Errorf("tls server quit with error: %v", err)

		}
	} else {
		if err := http.ListenAndServe(*listenAddress, nil); err != nil {

			log.WithFields(log.Fields{
				"err": err,
			}).Errorf("http server quit with error: %v", err)

		}
	}
}

// IndexHandler returns a http handler with the correct metricsPath
func IndexHandler(metricsPath string) http.HandlerFunc {
	indexHTML := `
<html>
	<head>
		<title>Beat Exporter</title>
	</head>
	<body>
		<h1>Beat Exporter</h1>
		<p>
			<a href='%s'>Metrics</a>
		</p>
	</body>
</html>
`
	index := []byte(fmt.Sprintf(strings.TrimSpace(indexHTML), metricsPath))

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(index)
	}
}

func loadBeatType(client *http.Client, url url.URL) (*collector.BeatInfo, error) {
	beatInfo := &collector.BeatInfo{}

	response, err := client.Get(url.String())
	if err != nil {
		return beatInfo, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Errorf("Beat URL: %q status code: %d", url.String(), response.StatusCode)
		return beatInfo, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Can't read body of response")
		return beatInfo, err
	}

	err = json.Unmarshal(bodyBytes, &beatInfo)
	if err != nil {
		log.Error("Could not parse JSON response for target")
		return beatInfo, err
	}

	log.WithFields(
		log.Fields{
			"beat":     beatInfo.Beat,
			"version":  beatInfo.Version,
			"name":     beatInfo.Name,
			"hostname": beatInfo.Hostname,
			"uuid":     beatInfo.UUID,
		}).Info("Target beat configuration loaded successfully!")

	return beatInfo, nil
}

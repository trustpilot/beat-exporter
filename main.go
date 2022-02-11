package main

import (
	"flag"
	"os"
	"time"

	"beat-exporter/internal/service"

	log "github.com/sirupsen/logrus"
)

var (
	listenAddress = flag.String("web.listen-address", ":9479", "Address to listen on for web interface and telemetry.")
	tlsCertFile   = flag.String("tls.certfile", "", "TLS certs file if you want to use tls instead of http")
	tlsKeyFile    = flag.String("tls.keyfile", "", "TLS key file if you want to use tls instead of http")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	beatURI       = flag.String("beat.uri", "http://localhost:5066", "HTTP API address of beat.")
	beatTimeout   = flag.Duration("beat.timeout", 10*time.Second, "Timeout for trying to get stats from beat.")
	showVersion   = flag.Bool("version", false, "Show version and exit")
	systemBeat    = flag.Bool("beat.system", false, "Expose system stats")
)

func init() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyMsg: "message",
		},
	})
}

func main() {
	params := service.NewParams(
		*listenAddress,
		*tlsCertFile,
		*tlsKeyFile,
		*metricsPath,
		*beatURI,
		*beatTimeout,
		*systemBeat,
	)

	if service.PrintVersion(*showVersion) {
		os.Exit(0)
	}

	exporter := service.NewExporter(params)
	httpHandler := service.NewHTTPHandler(params, exporter)
	go httpHandler.Run()

	log.WithFields(log.Fields{
		"addr": params.ListenAddress,
	}).Infof("Starting exporter with configured type: %s", exporter.BeatInfo.Beat)

	stopCh := make(chan bool)
	err := httpHandler.SetupServiceListener(stopCh, params.ServiceName, log.StandardLogger())
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("could not setup service listener: %v", err)
	}

	for {
		if <-stopCh {
			log.Info("Shutting down beats exporter")
			os.Exit(0) // signal received, stop gracefully
		}
	}
}

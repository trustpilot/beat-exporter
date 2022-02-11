package service

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type HTTPHandler struct {
	params   *Params
	registry *prometheus.Registry
}

func NewHTTPHandler(params *Params, exporter *Exporter) *HTTPHandler {
	h := &HTTPHandler{
		params:   params,
		registry: exporter.PrometheusRegistry,
	}

	http.HandleFunc("/", h.indexHandler())
	http.Handle(params.MetricsPath, h.prometheusHandler())
	return h
}

func (h *HTTPHandler) prometheusHandler() http.Handler {
	return promhttp.HandlerFor(
		h.registry,
		promhttp.HandlerOpts{
			ErrorLog:           log.New(),
			DisableCompression: false,
			ErrorHandling:      promhttp.ContinueOnError})
}

// IndexHandler returns a http handler with the correct metricsPath
func (h *HTTPHandler) indexHandler() http.HandlerFunc {

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

	index := []byte(fmt.Sprintf(strings.TrimSpace(indexHTML), h.params.MetricsPath))

	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(index)
	}
}

func (h *HTTPHandler) Run() {
	log.Info("Starting listener")
	params := h.params
	if params.TLSCertFile != "" && params.TLSKeyFile != "" {
		if err := http.ListenAndServeTLS(params.ListenAddress, params.TLSCertFile, params.TLSKeyFile, nil); err != nil {

			log.WithFields(log.Fields{
				"err": err,
			}).Errorf("tls server quit with error: %v", err)

		}
	} else {
		if err := http.ListenAndServe(params.ListenAddress, nil); err != nil {

			log.WithFields(log.Fields{
				"err": err,
			}).Errorf("http server quit with error: %v", err)

		}
	}
}

// SetupServiceListener setups signal handler
func (h *HTTPHandler) SetupServiceListener(stopCh chan<- bool, serviceName string, logger log.StdLogger) error {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
		logger.Printf("%s signal received: %v", serviceName, <-signals)
		stopCh <- true
		close(stopCh)
	}()
	return nil
}

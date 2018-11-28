package collector

import (
	"encoding/json"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"regexp"
	"time"
)

//Influxbeat json structure
type Influxbeat struct {
	Events struct {
		Active float64 `json:"active"`
		Added  float64 `json:"added"`
		Done   float64 `json:"done"`
	} `json:"events"`

	Harvester struct {
		Closed    float64 `json:"closed"`
		OpenFiles float64 `json:"open_files"`
		Running   float64 `json:"running"`
		Skipped   float64 `json:"skipped"`
		Started   float64 `json:"started"`
	} `json:"harvester"`

	Input struct {
		Log struct {
			Files struct {
				Renamed   float64 `json:"renamed"`
				Truncated float64 `json:"truncated"`
			} `json:"files"`
		} `json:"log"`
	} `json:"input"`
}

type influxbeatCollector struct {
	beatInfo     *BeatInfo
	stats        *Stats
	metrics      exportedMetrics
	influxClient influx.Client
}

// NewInluxbeatCollector constructor
func NewInluxbeatCollector(beatInfo *BeatInfo, stats *Stats) prometheus.Collector {
	influxClient, _ := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "",
		Password: "",
	})
	defer influxClient.Close()

	return &influxbeatCollector{
		beatInfo:     beatInfo,
		stats:        stats,
		influxClient: influxClient,
		metrics: exportedMetrics{
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "influxbeat", "events"),
					"influxbeat.events",
					nil, prometheus.Labels{"event": "active"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Events.Active },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "influxbeat", "events"),
					"influxbeat.events",
					nil, prometheus.Labels{"event": "added"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Events.Added },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "influxbeat", "events"),
					"influxbeat.events",
					nil, prometheus.Labels{"event": "done"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Events.Done },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "influxbeat", "harvester"),
					"influxbeat.harvester",
					nil, prometheus.Labels{"harvester": "closed"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Harvester.Closed },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "influxbeat", "harvester"),
					"influxbeat.harvester",
					nil, prometheus.Labels{"harvester": "open_files"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Harvester.OpenFiles },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "influxbeat", "harvester"),
					"influxbeat.harvester",
					nil, prometheus.Labels{"harvester": "running"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Harvester.Running },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "influxbeat", "harvester"),
					"influxbeat.harvester",
					nil, prometheus.Labels{"harvester": "skipped"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Harvester.Skipped },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "influxbeat", "harvester"),
					"influxbeat.harvester",
					nil, prometheus.Labels{"harvester": "started"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Harvester.Started },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "influxbeat", "input_log"),
					"influxbeat.input_log",
					nil, prometheus.Labels{"files": "renamed"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Input.Log.Files.Renamed },
				valType: prometheus.UntypedValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.Beat, "influxbeat", "input_log"),
					"influxbeat.input_log",
					nil, prometheus.Labels{"files": "truncated"},
				),
				eval:    func(stats *Stats) float64 { return stats.Filebeat.Input.Log.Files.Truncated },
				valType: prometheus.UntypedValue,
			},
		},
	}
}

// Describe returns all descriptions of the collector.
func (c *influxbeatCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.desc
	}
}

// Collect returns the current state of all metrics of the collector.
func (c *influxbeatCollector) Collect(ch chan<- prometheus.Metric) {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "dashbase_metric",
		Precision: "s",
	})
	if err != nil {
		log.Errorf("Failed create influx batch points: " + err.Error())
		return
	}

	for _, i := range c.metrics {
		ch <- prometheus.MustNewConstMetric(i.desc, i.valType, i.eval(c.stats))

		var eventsStats map[string]interface{}
		eventsBytes, _ := json.Marshal(c.stats.Filebeat.Events)
		json.Unmarshal(eventsBytes, &eventsStats)
		var harvesterStats map[string]interface{}
		harvesterBytes, _ := json.Marshal(c.stats.Filebeat.Harvester)
		json.Unmarshal(harvesterBytes, &harvesterStats)
		var inputStats map[string]interface{}
		inputBytes, _ := json.Marshal(c.stats.Filebeat.Input)
		json.Unmarshal(inputBytes, &inputStats)

		// TODO: delete regex, now prometheus not expose fqName public
		fqNameRegexp := regexp.MustCompile(`^Desc{fqName: "(\w+)", help:`)
		hits := fqNameRegexp.FindStringSubmatch(i.desc.String())
		fqName := ""
		if hits != nil && len(hits) >= 2 {
			fqName = hits[1]
		}

		eventsPt, _ := influx.NewPoint(
			fqName, map[string]string{}, eventsStats, time.Now(),
		)
		harvesterPt, _ := influx.NewPoint(
			fqName, map[string]string{}, harvesterStats, time.Now(),
		)
		inputPt, _ := influx.NewPoint(
			fqName, map[string]string{}, inputStats, time.Now(),
		)
		bp.AddPoint(eventsPt)
		bp.AddPoint(harvesterPt)
		bp.AddPoint(inputPt)
	}

	if err := c.influxClient.Write(bp); err != nil {
		log.Error(err)
	}
}

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
	"strings"
)

type apmServerCollector struct {
	beatInfo *BeatInfo
	stats    *Stats
	metrics  exportedMetrics
}

// NewApmServerCollector constructor
func NewApmServerCollector(beatInfo *BeatInfo, stats *Stats) prometheus.Collector {
	return &apmServerCollector{
		beatInfo: beatInfo,
		stats:    stats,
		metrics:  FlattenMetrics(beatInfo, stats.ApmServer, "", nil),
	}
}

// Describe returns all descriptions of the collector.
func (c *apmServerCollector) Describe(ch chan<- *prometheus.Desc) {

	for _, metric := range c.metrics {
		ch <- metric.desc
	}

}

// Collect returns the current state of all metrics of the collector.
func (c *apmServerCollector) Collect(ch chan<- prometheus.Metric) {

	for _, i := range c.metrics {
		ch <- prometheus.MustNewConstMetric(i.desc, i.valType, i.eval(c.stats))
	}

}

// Transform ApmServer struct into flat structure
func FlattenMetrics(beatInfo *BeatInfo, obj interface{}, prefix string, prefixArray []string) []exportedMetric {

	var exported []exportedMetric

	v := reflect.ValueOf(obj)
	currentType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			exported = append(exported, FlattenMetrics(beatInfo, v.Field(i).Interface(), prefix+formatMetricName(currentType.Field(i).Name)+"_", append(prefixArray, currentType.Field(i).Name))...)
		} else {
			var currentMetricName = currentType.Field(i).Name

			var metric = exportedMetric{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(beatInfo.FormattedBeat(), prefix[:len(prefix)-1], formatMetricName(currentMetricName)),
					prefix+currentType.Field(i).Name,
					nil, nil,
				),
				eval:    func(stats *Stats) float64 { return resolveStats(stats, append(prefixArray, currentMetricName)) },
				valType: prometheus.CounterValue,
			}
			exported = append(exported, metric)
		}

	}
	return exported
}

// Resolve stats by the array of prefix fields
func resolveStats(stats *Stats, prefixArray []string) float64 {
	v := reflect.ValueOf(stats.ApmServer)
	for _, prefix := range prefixArray {
		if v.Kind() == reflect.Struct {
			v = v.FieldByName(prefix)
		} else {
			return v.Interface().(float64)
		}
	}
	return v.Interface().(float64)
}

// Format the metric name - used to transform Struct names (Capitalized) into lowercase format
func formatMetricName(name string) string {
	return strings.ToLower(name)
}

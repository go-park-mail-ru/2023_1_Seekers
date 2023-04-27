package prometheus

import "github.com/prometheus/client_golang/prometheus"

type MetricsHttp struct {
	Errors  *prometheus.CounterVec
	Hits    *prometheus.Counter
	Timings *prometheus.HistogramVec
}

func NewMetricsHttp(name string) *MetricsHttp {
	metrics := &MetricsHttp{}

	return metrics
}

package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type MetricsHttp struct {
	Errors  *prometheus.CounterVec
	Hits    prometheus.Counter
	Timings *prometheus.HistogramVec
}

func NewMetricsHttpServer(name string) (*MetricsHttp, error) {
	metrics := &MetricsHttp{
		Hits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: name + "_hits",
			Help: "counts all hits for http server",
		}),
		Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: name + "_errors",
			Help: "counts responses with error from http server",
		}, []string{"code", "path", "method"}),
		Timings: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: name + "_timings",
			Help: "measures duration of http request",
		}, []string{"code", "path", "method"}),
	}

	if err := prometheus.Register(metrics.Hits); err != nil {
		return nil, err
	}

	if err := prometheus.Register(metrics.Errors); err != nil {
		return nil, err
	}

	if err := prometheus.Register(metrics.Timings); err != nil {
		return nil, err
	}

	return metrics, nil
}

func RunHttpMetricsServer(address string) error {
	//use separated ServeMux to prevent handling on the global Mux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	log.Info("starting metrics server...")

	return http.ListenAndServe(address, mux)
}

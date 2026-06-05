package httpapi

import (
	"github.com/prometheus/client_golang/prometheus"
	promrecorder "github.com/slok/go-http-metrics/metrics/prometheus"
)

type MetricsConfig struct{}

func (MetricsConfig) Build() (any, error) {
	reg := prometheus.NewRegistry()
	return promrecorder.NewRecorder(promrecorder.Config{Registry: reg}), nil
}

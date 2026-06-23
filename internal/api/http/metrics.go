package http

import (
	"github.com/prometheus/client_golang/prometheus"
	promrecorder "github.com/slok/go-http-metrics/metrics/prometheus"
)

// Metrics is the AppResources wire type for HTTP metrics recorder.
type Metrics struct{}

func (Metrics) NewResource() (any, error) {
	reg := prometheus.NewRegistry()
	return promrecorder.NewRecorder(promrecorder.Config{Registry: reg}), nil
}

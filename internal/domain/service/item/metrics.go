package item

import (
	"errors"

	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/ops/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type serviceMetrics struct {
	ops *prometheus.CounterVec
}

func (s *Service) RegisterMetrics(reg *prometheus.Registry) error {
	s.metrics = &serviceMetrics{
		ops: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "demo_item_ops_total",
			Help: "Item domain operations by operation and result.",
		}, []string{"operation", "result"}),
	}
	return reg.Register(s.metrics.ops)
}

func (o *observeService) RegisterMetrics(reg *prometheus.Registry) error {
	return o.Service.RegisterMetrics(reg)
}

func (s *Service) recordOp(operation, result string) {
	if s.metrics == nil || s.metrics.ops == nil {
		return
	}
	s.metrics.ops.WithLabelValues(operation, result).Inc()
}

func classifyErr(err error) string {
	if err == nil {
		return "ok"
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return "not_found"
	case errors.Is(err, domain.ErrInvalidInput):
		return "invalid"
	default:
		return "error"
	}
}

var (
	_ metrics.MetricsContributor = (*Service)(nil)
	_ metrics.MetricsContributor = (*observeService)(nil)
)

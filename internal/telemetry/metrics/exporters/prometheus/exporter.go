package prometheus

import (
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func New() (metric.Reader, error) {
	return prometheus.New()
}

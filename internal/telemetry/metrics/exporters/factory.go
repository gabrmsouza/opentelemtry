package exporters

import (
	"context"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/exporters/prometheus"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/exporters/stdout"
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Reader interface {
	Start(ctx context.Context) (metric.Reader, error)
	Shutdown(context.Context) error
}

func New(exporter properties.MetricExporter) Reader {
	switch exporter.Type {
	case properties.PrometheusReader:
		return prometheus.New()
	case properties.StdoutReader:
		return stdout.New()
	default:
		return stdout.New()
	}
}

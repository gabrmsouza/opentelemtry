package exporters

import (
	"context"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/exporters/prometheus"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/exporters/stdout"
	"go.opentelemetry.io/otel/sdk/metric"
)

func New(ctx context.Context, exporter string) (metric.Reader, error) {
	switch exporter {
	case "prometheus":
		return prometheus.New()
	case "stdout":
		return stdout.New()
	default:
		return stdout.New()
	}
}

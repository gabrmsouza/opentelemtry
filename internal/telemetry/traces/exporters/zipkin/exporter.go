package zipkin

import (
	"context"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/properties"
	"go.opentelemetry.io/otel/exporters/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type zipkinExporter struct {
	props properties.Exporter
}

func New(props properties.Exporter) *zipkinExporter {
	return &zipkinExporter{props: props}
}

func (h *zipkinExporter) GetExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return zipkin.New(h.props.EndpointURL)
}

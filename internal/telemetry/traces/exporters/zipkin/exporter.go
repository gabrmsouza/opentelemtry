package zipkin

import (
	"context"

	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"
	"go.opentelemetry.io/otel/exporters/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Exporter struct {
	props    properties.TraceExporter
	shutdown func(ctx context.Context) error
}

func New(props properties.TraceExporter) *Exporter {
	return &Exporter{props: props}
}

func (h *Exporter) Start(_ context.Context) (sdktrace.SpanExporter, error) {
	exporter, err := zipkin.New(h.props.EndpointURL)
	if err != nil {
		return nil, err
	}
	h.shutdown = exporter.Shutdown
	return exporter, nil
}

func (h Exporter) Shutdown(ctx context.Context) error {
	return h.shutdown(ctx)
}

package grpc

import (
	"context"

	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Exporter struct {
	props    properties.TraceExporter
	shutdown func(ctx context.Context) error
}

func New(props properties.TraceExporter) *Exporter {
	return &Exporter{props: props}
}

func (h *Exporter) Start(ctx context.Context) (sdktrace.SpanExporter, error) {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpointURL(h.props.EndpointURL),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	h.shutdown = exporter.Shutdown
	return exporter, nil
}

func (h Exporter) Shutdown(ctx context.Context) error {
	return h.shutdown(ctx)
}

package grpc

import (
	"context"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/properties"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type grpcExporter struct {
	props properties.Exporter
}

func New(props properties.Exporter) *grpcExporter {
	return &grpcExporter{props: props}
}

func (g *grpcExporter) GetExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpointURL(g.props.EndpointURL),
		otlptracegrpc.WithInsecure(),
	)
}

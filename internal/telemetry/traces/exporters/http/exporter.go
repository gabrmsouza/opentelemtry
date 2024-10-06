package http

import (
	"context"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/properties"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type httpExporter struct {
	props properties.Exporter
}

func New(props properties.Exporter) *httpExporter {
	return &httpExporter{props: props}
}

func (h *httpExporter) GetExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	e, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpointURL(h.props.EndpointURL),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	return e, nil
}

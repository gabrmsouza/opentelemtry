package stdout

import (
	"context"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type stdoutExporter struct {
}

func New() *stdoutExporter {
	return &stdoutExporter{}
}

func (s stdoutExporter) GetExporter(_ context.Context) (sdktrace.SpanExporter, error) {
	exporter, err := stdouttrace.New()
	return exporter, err
}

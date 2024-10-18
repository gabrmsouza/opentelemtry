package stdout

import (
	"context"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Exporter struct {
	shutdown func(ctx context.Context) error
}

func New() *Exporter {
	return &Exporter{}
}

func (s *Exporter) Start(_ context.Context) (sdktrace.SpanExporter, error) {
	exporter, err := stdouttrace.New()
	if err != nil {
		return nil, err
	}
	s.shutdown = exporter.Shutdown
	return exporter, nil
}

func (s Exporter) Shutdown(ctx context.Context) error {
	return s.shutdown(ctx)
}

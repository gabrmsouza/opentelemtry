package stdout

import (
	"context"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Reader struct {
	shutdown func(ctx context.Context) error
}

func New() *Reader {
	return &Reader{}
}

func (s *Reader) Start(_ context.Context) (metric.Reader, error) {
	exporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}
	s.shutdown = exporter.Shutdown
	return metric.NewPeriodicReader(exporter), nil
}

func (s Reader) Shutdown(ctx context.Context) error {
	return s.shutdown(ctx)
}

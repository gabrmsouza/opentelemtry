package prometheus

import (
	"context"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Reader struct {
	shutdown func(ctx context.Context) error
}

func New() *Reader {
	return &Reader{}
}

func (p *Reader) Start(_ context.Context) (metric.Reader, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	p.shutdown = exporter.Shutdown
	return exporter, nil
}

func (p Reader) Shutdown(ctx context.Context) error {
	return p.shutdown(ctx)
}

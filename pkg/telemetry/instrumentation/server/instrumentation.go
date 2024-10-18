package instrumentation

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	HttpServerRequestDurationMetricName = "http.server.request.duration"

	ApplicationNameAttribute        = "application"
	HttpRouteAttribute              = "http.route"
	HttpRequestMethodAttribute      = "http.request.method"
	HttpResponseStatusCodeAttribute = "http.response.status_code"
	ServerAddressAttribute          = "server.address"
)

type instrumentation struct {
	serverLatencyMeasure metric.Float64Histogram
}

func New(serviceName string) (*instrumentation, error) {
	meter := otel.Meter(serviceName)
	serverLatencyMeasure, err := meter.Float64Histogram(
		HttpServerRequestDurationMetricName,
		metric.WithDescription("Duration of HTTP server requests."),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10),
	)
	if err != nil {
		return nil, err
	}
	return &instrumentation{
		serverLatencyMeasure: serverLatencyMeasure,
	}, nil
}

func (m *instrumentation) RecordServerLatency(ctx context.Context, incr float64, options ...metric.RecordOption) {
	m.serverLatencyMeasure.Record(ctx, incr, options...)
}

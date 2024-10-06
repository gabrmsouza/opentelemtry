package exporters

import (
	"context"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/properties"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters/grpc"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters/http"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters/stdout"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters/zipkin"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Exporter interface {
	GetExporter(ctx context.Context) (sdktrace.SpanExporter, error)
}

func New(props properties.Exporter) Exporter {
	switch props.Type {
	case "http":
		return http.New(props)
	case "grpc":
		return grpc.New(props)
	case "stdout":
		return stdout.New()
	case "zipkin":
		return zipkin.New(props)
	default:
		return stdout.New()
	}
}

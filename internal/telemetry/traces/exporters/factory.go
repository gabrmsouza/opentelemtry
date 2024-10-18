package exporters

import (
	"context"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters/grpc"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters/http"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters/stdout"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters/zipkin"
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Exporter interface {
	Start(ctx context.Context) (sdktrace.SpanExporter, error)
	Shutdown(ctx context.Context) error
}

func New(props properties.TraceExporter) Exporter {
	switch props.Type {
	case properties.HttpExporter:
		return http.New(props)
	case properties.GrpcExporter:
		return grpc.New(props)
	case properties.StdoutExporter:
		return stdout.New()
	case properties.ZipkinExporter:
		return zipkin.New(props)
	default:
		return stdout.New()
	}
}

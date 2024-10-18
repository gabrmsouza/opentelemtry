package resource

import (
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func New(props properties.Service) (*resource.Resource, error) {
	return resource.Merge(
		resource.Environment(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(props.Name),
			semconv.ServiceVersion(props.Version),
			semconv.DeploymentEnvironment("local"),
		),
	)
}

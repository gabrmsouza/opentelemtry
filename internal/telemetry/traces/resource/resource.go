package resource

import (
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/properties"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func New(props properties.Service) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(props.Name),
			semconv.ServiceVersion(props.Version),
			semconv.DeploymentEnvironment("local"),
		),
	)
}

package provider

import (
	"context"
	"fmt"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/exporters"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/properties"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/resource"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Provider struct {
	props properties.Telemetry
}

func New(props properties.Telemetry) *Provider {
	return &Provider{props}
}

func (p *Provider) Start(ctx context.Context) func(ctx context.Context) error {
	if !p.props.Enabled {
		return func(ctx context.Context) error {
			return nil
		}
	}
	exporter, err := exporters.New(ctx, p.props.Metric.Exporter.Type)
	if err != nil {
		panic(err)
	}

	r, err := resource.New(p.props.Service)
	if err != nil {
		panic(err)
	}

	if p.props.Trace.Enabled {
		otel.SetMeterProvider(sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(exporter),
			sdkmetric.WithResource(r)),
		)
	} else {
		otel.SetMeterProvider(noop.NewMeterProvider())
	}
	return exporter.Shutdown
}

func (p *Provider) GetMeter() metric.Meter {
	return otel.Meter(fmt.Sprintf("io.opentelemetry.metrics.%s", p.props.Service.Name))
}

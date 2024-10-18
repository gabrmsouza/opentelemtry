package provider

import (
	"context"
	"fmt"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/exporters"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/resource"
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Provider struct {
	props  properties.Telemetry
	reader exporters.Reader
}

func New(props properties.Telemetry, reader exporters.Reader) *Provider {
	return &Provider{
		props:  props,
		reader: reader,
	}
}

func (p *Provider) Start(ctx context.Context) error {
	reader, err := p.reader.Start(ctx)
	if err != nil {
		return err
	}

	r, err := resource.New(p.props.Service)
	if err != nil {
		return err
	}

	if !p.props.Metric.Enabled {
		otel.SetMeterProvider(noop.NewMeterProvider())
	} else {
		otel.SetMeterProvider(sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(reader),
			sdkmetric.WithResource(r)),
		)
	}
	return nil
}

func (p *Provider) Shutdown(ctx context.Context) error {
	return p.reader.Shutdown(ctx)
}

func (p *Provider) GetMeter() metric.Meter {
	return otel.Meter(fmt.Sprintf("io.opentelemetry.metrics.%s", p.props.Service.Name))
}

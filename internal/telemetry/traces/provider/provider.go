package provider

import (
	"context"
	"fmt"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/resource"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters"
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Provider struct {
	props    properties.Telemetry
	exporter exporters.Exporter
}

func New(props properties.Telemetry, exporter exporters.Exporter) *Provider {
	return &Provider{
		props:    props,
		exporter: exporter,
	}
}

func (p *Provider) Start(ctx context.Context) error {
	e, err := p.exporter.Start(ctx)
	if err != nil {
		return err
	}

	r, err := resource.New(p.props.Service)
	if err != nil {
		return err
	}

	if !p.props.Trace.Enabled {
		otel.SetTracerProvider(noop.NewTracerProvider())
	} else {
		otel.SetTracerProvider(sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(e)),
			sdktrace.WithResource(r),
		))
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return nil
}

func (p *Provider) GetTracer() trace.Tracer {
	return otel.Tracer(fmt.Sprintf("io.opentelemetry.traces.%s", p.props.Service.Name))
}

func (p *Provider) Shutdown(ctx context.Context) error {
	return p.exporter.Shutdown(ctx)
}

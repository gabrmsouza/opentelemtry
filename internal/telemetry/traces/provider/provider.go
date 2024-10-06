package provider

import (
	"context"
	"fmt"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/properties"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/resource"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
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
	exporter, err := exporters.New(p.props.Trace.Exporter).GetExporter(ctx)
	if err != nil {
		panic(err)
	}

	r, err := resource.New(p.props.Service)
	if err != nil {
		panic(err)
	}

	if p.props.Trace.Enabled {
		otel.SetTracerProvider(sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
			sdktrace.WithResource(r),
		))
	} else {
		otel.SetTracerProvider(noop.NewTracerProvider())
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return exporter.Shutdown
}

func (p *Provider) GetTracer() trace.Tracer {
	return otel.Tracer(fmt.Sprintf("io.opentelemetry.traces.%s", p.props.Service.Name))
}

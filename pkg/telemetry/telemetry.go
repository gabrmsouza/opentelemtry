package telemetry

import (
	"context"
	"fmt"
	metricexporter "github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/exporters"
	metricprovider "github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/provider"
	traceexporter "github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/exporters"
	traceprovider "github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/provider"
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type telemetry struct {
	traceProvider  *traceprovider.Provider
	metricProvider *metricprovider.Provider
}

func New(props properties.Telemetry) *telemetry {
	return &telemetry{
		traceProvider:  traceprovider.New(props, traceexporter.New(props.Trace.Exporter)),
		metricProvider: metricprovider.New(props, metricexporter.New(props.Metric.Exporter)),
	}
}

func (t telemetry) StartTraceProvider(ctx context.Context) {
	if err := t.traceProvider.Start(ctx); err != nil {
		log.Println(fmt.Sprintf("failed to starting trace provider [err:%s]", err.Error()))
	}
}

func (t telemetry) ShutdownTraceProvider(ctx context.Context) {
	if err := t.traceProvider.Shutdown(ctx); err != nil {
		log.Println(fmt.Sprintf("failed to shutdown trace provider [err:%s]", err.Error()))
	}
}

func (t telemetry) GetTracer() trace.Tracer {
	return t.traceProvider.GetTracer()
}

func (t telemetry) StartMetricProvider(ctx context.Context) {
	if err := t.metricProvider.Start(ctx); err != nil {
		log.Println(fmt.Sprintf("failed to starting metric provider [err:%s]", err.Error()))
	}
}

func (t telemetry) ShutdownMetricProvider(ctx context.Context) {
	if err := t.metricProvider.Shutdown(ctx); err != nil {
		log.Println(fmt.Sprintf("failed to shutdown metric provider [err:%s]", err.Error()))
	}
}

func (t telemetry) GetMeter() metric.Meter {
	return t.metricProvider.GetMeter()
}

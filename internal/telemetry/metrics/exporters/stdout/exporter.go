package stdout

import (
    "go.opentelemetry.io/otel/sdk/metric"
    "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
)

func New() (metric.Reader, error) {
    exporter, err := stdoutmetric.New()
    if err != nil {
        return nil, err
    }
    return metric.NewPeriodicReader(exporter), nil
}
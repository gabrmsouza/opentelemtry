package properties

type MetricReaderType string
type TraceExporterType string

const (
    // Traces exporters
    HttpExporter   TraceExporterType = "http"
    GrpcExporter   TraceExporterType = "grpc"
    StdoutExporter TraceExporterType = "stdout"
    ZipkinExporter TraceExporterType = "zipkin"

    // Metrics readers
    PrometheusReader MetricReaderType = "prometheus"
    StdoutReader     MetricReaderType = "stdout"
)

type (
    Telemetry struct {
       Service Service `yaml:"service"`
       Trace   Trace   `yaml:"traces"`
       Metric  Metric  `yaml:"metrics"`
    }

    Service struct {
       Name    string `yaml:"name"`
       Version string `yaml:"version"`
    }

    Metric struct {
       Enabled  bool           `yaml:"enabled"`
       Exporter MetricExporter `yaml:"exporter"`
    }

    MetricExporter struct {
       Type     MetricReaderType `yaml:"type"`
       Endpoint string           `yaml:"endpoint"`
       Port     string           `yaml:"port"`
    }

    Trace struct {
       Enabled  bool          `yaml:"enabled"`
       Exporter TraceExporter `yaml:"exporter"`
    }

    TraceExporter struct {
       Type        TraceExporterType `yaml:"type"`
       EndpointURL string            `yaml:"endpoint-url"`
    }
)
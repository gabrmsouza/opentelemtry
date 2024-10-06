package properties

// Root props
type (
	Root struct {
		Telemetry Telemetry `yaml:"telemetry"`
	}

	Telemetry struct {
		Enabled bool    `yaml:"enabled"`
		Service Service `yaml:"service"`
		Trace   Trace   `yaml:"traces"`
		Metric  Metric  `yaml:"metrics"`
	}
)

func (t Telemetry) GetExporter() Exporter {
	return t.Trace.Exporter
}

// Service props
type Service struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// Trace props
type (
	Exporter struct {
		Type        string `yaml:"type"`
		EndpointURL string `yaml:"endpoint-url"`
		Insecure    bool   `yaml:"insecure"`
	}

	Trace struct {
		Enabled  bool     `yaml:"enabled"`
		Exporter Exporter `yaml:"exporter"`
	}
)

// Metrics props
type (
	Metric struct {
		Enabled  bool           `yaml:"enabled"`
		Exporter MetricExporter `yaml:"exporter"`
	}

	MetricExporter struct {
		Type     string `yaml:"type"`
		Endpoint string `yaml:"endpoint"`
		Port     string `yaml:"port"`
	}
)

package prometheus

import (
	"net/http"

	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricServer struct {
	endpoint string
	port     string
}

func New(props properties.MetricExporter) *MetricServer {
	return &MetricServer{
		endpoint: props.Endpoint,
		port:     props.Port,
	}
}

func (s *MetricServer) Serve() error {
	prometheus := promhttp.Handler()
	router := mux.NewRouter()
	router.HandleFunc(s.endpoint, func(w http.ResponseWriter, r *http.Request) {
		prometheus.ServeHTTP(w, r)
	})
	err := http.ListenAndServe(":"+s.port, router)
	return err
}

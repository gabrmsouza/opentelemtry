package prometheus

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricServer struct {
	endpoint string
	port     string
}

func New(endpoint, port string) *MetricServer {
	return &MetricServer{
		endpoint: endpoint,
		port:     port,
	}
}

func (s *MetricServer) Serve(errCh chan error) {
	prometheus := promhttp.Handler()
	router := mux.NewRouter()
	router.HandleFunc(s.endpoint, func(w http.ResponseWriter, r *http.Request) {
		prometheus.ServeHTTP(w, r)
	})

	errCh <- http.ListenAndServe(":"+s.port, router)
}

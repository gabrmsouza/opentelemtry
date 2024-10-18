package metrics

import (
	"log"
	"net/http"
	"strconv"
	"time"

	instrumentation "github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/instrumentation/server"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/metric"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Middleware(serviceName string) mux.MiddlewareFunc {
	ins, err := instrumentation.New(serviceName)
	if err != nil {
		log.Fatalf("failed to create instrumentation [err:%s]", err.Error())
	}
	app := attribute.String(instrumentation.ApplicationNameAttribute, serviceName)
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			l := wrapResponseWriter(w)
			handler.ServeHTTP(l, r)
			elapsedTime := float64(time.Since(startTime)) / float64(time.Second)
			server := attribute.String(instrumentation.ServerAddressAttribute, r.Host)
			path := attribute.String(instrumentation.HttpRouteAttribute, r.URL.Path)
			method := attribute.String(instrumentation.HttpRequestMethodAttribute, r.Method)
			status := attribute.String(instrumentation.HttpResponseStatusCodeAttribute, strconv.Itoa(l.statusCode))
			attributes := sdkmetric.WithAttributes(app, server, path, method, status, server)
			ins.RecordServerLatency(r.Context(), elapsedTime, attributes)
		})
	}
}

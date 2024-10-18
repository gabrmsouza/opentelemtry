package tracing

import (
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
)

var endpointsToExcludeFromTrace = []string{"/health", "/actuator/health", "/metrics"}

func Middleware(serviceName string) mux.MiddlewareFunc {
	return otelmux.Middleware(serviceName,
		otelmux.WithTracerProvider(otel.GetTracerProvider()),
		otelmux.WithPropagators(otel.GetTextMapPropagator()),
		otelmux.WithFilter(func(request *http.Request) bool {
			return slices.Index(endpointsToExcludeFromTrace, request.URL.Path) == -1
		}),
	)
}

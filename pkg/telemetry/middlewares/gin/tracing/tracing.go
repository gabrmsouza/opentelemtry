package tracing

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

var endpointsToExcludeFromTrace = []string{"/health", "/actuator/health", "/metrics"}

func Middleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName,
		otelgin.WithTracerProvider(otel.GetTracerProvider()),
		otelgin.WithPropagators(otel.GetTextMapPropagator()),
		otelgin.WithFilter(func(request *http.Request) bool {
			return slices.Index(endpointsToExcludeFromTrace, request.URL.Path) == -1
		}),
	)
}

package metrics

import (
	"log"
	"strconv"
	"time"

	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/instrumentation/server"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func Middleware(serviceName string) gin.HandlerFunc {
	ins, err := instrumentation.New(serviceName)
	if err != nil {
		log.Fatalf("failed to create instrumentation [err:%s]", err.Error())
	}
	app := attribute.String(instrumentation.ApplicationNameAttribute, serviceName)
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		elapsedTime := float64(time.Since(startTime)) / float64(time.Second)
		path := attribute.String(instrumentation.HttpRouteAttribute, c.Request.URL.Path)
		method := attribute.String(instrumentation.HttpRequestMethodAttribute, c.Request.Method)
		status := attribute.String(instrumentation.HttpResponseStatusCodeAttribute, strconv.Itoa(c.Writer.Status()))
		server := attribute.String(instrumentation.ServerAddressAttribute, c.Request.Host)
		attributes := metric.WithAttributes(app, server, path, method, status)
		ins.RecordServerLatency(c.Request.Context(), elapsedTime, attributes)
	}
}

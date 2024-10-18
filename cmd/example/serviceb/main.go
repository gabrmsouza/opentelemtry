package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry"
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/middlewares/mux/metrics"
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type Root struct {
	Telemetry properties.Telemetry `yaml:"telemetry"`
}

func main() {
	// os.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.name=serviceb")
	data, err := os.ReadFile("../../../properties.yml")
	if err != nil {
		panic(err)
	}

	var root Root
	if err := yaml.Unmarshal(data, &root); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	root.Telemetry.Service.Name = "serviceb"
	tel := telemetry.New(root.Telemetry)
	tel.StartTraceProvider(ctx)
	defer tel.ShutdownTraceProvider(ctx)

	tracer := tel.GetTracer()

	router := mux.NewRouter()
	router.Use(metrics.Middleware(root.Telemetry.Service.Name))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(ctx, "processing-request")
		w.WriteHeader(200)
		w.Write([]byte(`{"message":"hello world"}`))
		span.End()
	})
	fmt.Println("server linten on port 8081")
	http.ListenAndServe(":8081", router)
}

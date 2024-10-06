package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	traceprovider "github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/provider"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/properties"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"
)

var root properties.Root
var tracer trace.Tracer

func main() {
	data, err := os.ReadFile("../../properties.yml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &root); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	root.Telemetry.Service.Name = "service-b"

	p := traceprovider.New(root.Telemetry)
	tracer = p.GetTracer()
	shutdown := p.Start(ctx)
	defer func() {
		if err = shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	router := mux.NewRouter()
	router.Use(otelmux.Middleware(root.Telemetry.Service.Name))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(ctx, "processing-request")
		w.WriteHeader(200)
		w.Write([]byte(`{"message":"hello world"}`))
		span.End()
	})
	fmt.Println("server linten on port 8081")
	http.ListenAndServe(":8081", router)
}

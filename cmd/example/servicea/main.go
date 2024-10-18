package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/server/prometheus"
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry"
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/middlewares/mux/metrics"
	"github.com/gabrmsouza/fullcycle/opentelemetry/pkg/telemetry/properties"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"
)

type Root struct {
	Telemetry properties.Telemetry `yaml:"telemetry"`
}

var tracer trace.Tracer

func main() {
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

	tel := telemetry.New(root.Telemetry)
	tel.StartTraceProvider(ctx)
	defer tel.ShutdownTraceProvider(ctx)

	tracer = tel.GetTracer()

	tel.StartMetricProvider(ctx)
	defer tel.ShutdownMetricProvider(ctx)

	go func() {
		fmt.Println("metrics server linten on port 3001")
		srv := prometheus.New(root.Telemetry.Metric.Exporter)
		if err := srv.Serve(); err != nil {
			fmt.Println(err)
		}
	}()

	router := mux.NewRouter()
	router.Use(metrics.Middleware(root.Telemetry.Service.Name))
	router.HandleFunc("/", homeHandler)
	fmt.Println("server lintem on port 8080")
	http.ListenAndServe(":8080", router)
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	// rotina 1 - Process File
	ctx, processFile := tracer.Start(ctx, "process-file")
	time.Sleep(time.Millisecond * 100)
	processFile.End()

	// rotina 2 - Fazer Request HTTP
	ctx, httpCall := tracer.Start(ctx, "request-remote-json")
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/", nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req) // chamo a requisição
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Millisecond * 300)
	httpCall.End()

	// rotina 3 - Exibir resultado
	_, renderContent := tracer.Start(ctx, "render-content")
	time.Sleep(time.Millisecond * 100)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(body))
	renderContent.End()
}

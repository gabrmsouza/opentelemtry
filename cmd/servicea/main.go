package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	metricprovider "github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/provider"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/metrics/server/prometheus"
	"github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/properties"
	traceprovider "github.com/gabrmsouza/fullcycle/opentelemetry/internal/telemetry/traces/provider"

	"math/rand"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/metric"
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

	p := traceprovider.New(root.Telemetry)
	tracer = p.GetTracer()
	shutdown := p.Start(ctx)
	defer func() {
		if err = shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	mp := metricprovider.New(root.Telemetry)
	mpshutdown := mp.Start(ctx)
	meter := mp.GetMeter()
	defer func() {
		if err = mpshutdown(ctx); err != nil {
			panic(err)
		}
	}()

	errCh := make(chan error, 1)
	go func() {
		fmt.Println("metrics server lintem on port 3001")
		srv := prometheus.New(root.Telemetry.Metric.Exporter.Endpoint, root.Telemetry.Metric.Exporter.Port)
		srv.Serve(errCh)
		if err := <-errCh; err != nil {
			panic(err)
		}
	}()

	opt := metric.WithAttributes(
		attribute.Key("application").String(root.Telemetry.Service.Name),
	)

	// This is the equivalent of prometheus.NewCounterVec
	counter, err := meter.Float64Counter(
		root.Telemetry.Service.Name,
		metric.WithDescription("a simple counter"),
	)
	if err != nil {
		panic(err)
	}
	counter.Add(ctx, 5, opt)

	gauge, err := meter.Float64ObservableGauge(
		root.Telemetry.Service.Name,
		metric.WithDescription("a fun little gauge"),
	)
	if err != nil {
		panic(err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	_, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		n := -10. + rng.Float64()*(90.) // [-10, 100)
		o.ObserveFloat64(gauge, n, opt)
		return nil
	}, gauge)
	if err != nil {
		log.Fatal(err)
	}

	// This is the equivalent of prometheus.NewHistogramVec
	histogram, err := meter.Float64Histogram(
		root.Telemetry.Service.Name,
		metric.WithDescription("a histogram with custom buckets and rename"),
		metric.WithExplicitBucketBoundaries(64, 128, 256, 512, 1024, 2048, 4096),
	)
	if err != nil {
		panic(err)
	}
	histogram.Record(ctx, 136, opt)
	histogram.Record(ctx, 64, opt)
	histogram.Record(ctx, 701, opt)
	histogram.Record(ctx, 830, opt)

	router := mux.NewRouter()
	router.Use(otelmux.Middleware(root.Telemetry.Service.Name))
	router.HandleFunc("/", homeHandler)
	fmt.Println("server lintem on port 8080")
	http.ListenAndServe(":8080", router)
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := baggage.ContextWithoutBaggage(request.Context())

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

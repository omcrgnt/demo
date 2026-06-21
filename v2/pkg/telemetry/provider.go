package telemetry

import (
	"context"
	"fmt"
	"net"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Provider configures the process TracerProvider and lives in res.
type Provider struct {
	tp *sdktrace.TracerProvider
}

func (p *Provider) Deps() []any {
	return nil
}

func (p *Provider) Inject(_ []any) {
	if p.tp == nil {
		return
	}
	otel.SetTracerProvider(p.tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

// Close shuts down the TracerProvider (runner.Closer).
func (p *Provider) Close(ctx context.Context) error {
	if p.tp == nil {
		return nil
	}
	return p.tp.Shutdown(ctx)
}

func newProvider(host string, port uint32, insecure bool, serviceName string) (*Provider, error) {
	endpoint := formatEndpoint(host, port)

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
	}
	if insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("telemetry: otlp http exporter: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("telemetry: resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	return &Provider{tp: tp}, nil
}

func formatEndpoint(host string, port uint32) string {
	return net.JoinHostPort(host, fmt.Sprintf("%d", port))
}

func noopTracerProvider() *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider()
}

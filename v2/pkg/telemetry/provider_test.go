package telemetry

import (
	"context"
	"testing"

	commonv1 "github.com/omcrgnt/proto/gen/go/common/v1"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func newTestProvider(t *testing.T) (*Provider, *tracetest.InMemoryExporter) {
	t.Helper()

	exp := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exp))
	return &Provider{tp: tp}, exp
}

func TestProvider_InjectRecordsSpan(t *testing.T) {
	p, exp := newTestProvider(t)
	p.Inject(nil)

	_, span := otel.Tracer("test").Start(context.Background(), "op")
	span.End()

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("spans: got %d want 1", len(spans))
	}
	if spans[0].Name != "op" {
		t.Fatalf("span name: got %q", spans[0].Name)
	}

	if err := p.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestConfig_Build(t *testing.T) {
	raw, err := (Config{
		ServiceName: &commonv1.Label{Value: "demo"},
		Host:        &commonv1.Host{Value: "127.0.0.1"},
		Port:        &commonv1.Port{Value: 4318},
		Insecure:    true,
	}).Build()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := raw.(*Provider); !ok {
		t.Fatalf("got %T", raw)
	}
}

func TestFormatEndpoint(t *testing.T) {
	if got := formatEndpoint("localhost", 4318); got != "localhost:4318" {
		t.Fatalf("got %q", got)
	}
}

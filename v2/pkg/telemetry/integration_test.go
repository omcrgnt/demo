package telemetry_test

import (
	"context"
	"testing"

	commonv1 "github.com/omcrgnt/proto/gen/go/common/v1"
	"github.com/omcrgnt/demo/v2/pkg/builder"
	"github.com/omcrgnt/demo/v2/pkg/res"
	"github.com/omcrgnt/demo/v2/pkg/res/restest"
	"github.com/omcrgnt/demo/v2/pkg/sdi"
	"github.com/omcrgnt/demo/v2/pkg/telemetry"
	"go.opentelemetry.io/otel"
)

func setupUseDefaults(t *testing.T) {
	t.Helper()
	restest.ResetGlobal()
	_ = res.AddToGlobalWithTags(telemetry.DefaultTraceConfig(), res.TagReplaceable)
	if err := builder.Build(res.Global()); err != nil {
		t.Fatal(err)
	}
}

func TestIntegration_resolveWithUseDefault(t *testing.T) {
	setupUseDefaults(t)

	if err := sdi.Resolve(res.Global()); err != nil {
		t.Fatal(err)
	}

	_, span := otel.Tracer("test").Start(context.Background(), "resolve-default")
	defer span.End()

	if !span.SpanContext().TraceID().IsValid() {
		t.Fatal("expected valid trace id after resolve")
	}
}

func TestIntegration_dedupUserOverride(t *testing.T) {
	setupUseDefaults(t)

	if err := res.Global().Add(telemetry.Config{
		ServiceName: &commonv1.Label{Value: "demo"},
		Host:        &commonv1.Host{Value: "127.0.0.1"},
		Port:        &commonv1.Port{Value: 4318},
		Insecure:    true,
	}); err != nil { //nolint:forbidigo // simulates ecfg.Register
		t.Fatal(err)
	}
	if err := builder.Build(res.Global()); err != nil {
		t.Fatal(err)
	}

	if err := sdi.Resolve(res.Global()); err != nil {
		t.Fatal(err)
	}

	_, span := otel.Tracer("test").Start(context.Background(), "user-override")
	defer span.End()

	if !span.SpanContext().TraceID().IsValid() {
		t.Fatal("expected valid trace id after user override resolve")
	}
}

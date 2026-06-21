package main

import (
	"testing"

	"github.com/omcrgnt/demo/v2/internal/api/http"
	orderhttp "github.com/omcrgnt/demo/v2/internal/api/http/order"
	"github.com/omcrgnt/demo/v2/internal/domain/service/item"
	"github.com/omcrgnt/demo/v2/pkg/app"
	"github.com/omcrgnt/demo/v2/pkg/builder"
	"github.com/omcrgnt/demo/v2/pkg/ecfg"
	_ "github.com/omcrgnt/demo/v2/pkg/logger/use"
	"github.com/omcrgnt/demo/v2/pkg/res"
	srvhttp "github.com/omcrgnt/demo/v2/pkg/srv-http"
)

func TestAppResources_Apply(t *testing.T) {
	t.Setenv("DEMO_APP_SHUTDOWN_TIMEOUT", "5s")
	t.Setenv("DEMO_SERVER_HTTP_ITEM_LABEL", "demo")
	t.Setenv("DEMO_SERVER_HTTP_ITEM_HOST", "127.0.0.1")
	t.Setenv("DEMO_SERVER_HTTP_ITEM_PORT", "8080")
	t.Setenv("DEMO_SERVICE_ITEM_MAX_LIST_LEN", "100")
	t.Setenv("DEMO_SERVER_HTTP_ORDER_LABEL", "orders")
	t.Setenv("DEMO_SERVER_HTTP_ORDER_HOST", "127.0.0.1")
	t.Setenv("DEMO_SERVER_HTTP_ORDER_PORT", "8081")

	var ar AppResources

	reg := res.New()
	if err := builder.Seed(reg, &ar); err != nil {
		t.Fatal(err)
	}
	if err := ecfg.Apply(reg, &ar, ecfg.WithPrefix(envPrefix)); err != nil {
		t.Fatal(err)
	}

	m, ok := builder.SeedMapFor(reg)
	if !ok {
		t.Fatal("expected seed map")
	}
	spec, ok := m["ServerHTTPItem"]
	if !ok {
		t.Fatal("expected ServerHTTPItem spec")
	}
	cfg, ok := spec.(*srvhttp.Config[*http.API])
	if !ok {
		t.Fatalf("unexpected spec type %T", spec)
	}
	if cfg.Label.GetValue() != "demo" {
		t.Fatalf("label: got %q", cfg.Label.GetValue())
	}

	spec, ok = m["ServiceItem"]
	if !ok {
		t.Fatal("expected ServiceItem spec")
	}
	itemSpec, ok := spec.(*item.Spec)
	if !ok {
		t.Fatalf("unexpected item spec type %T", spec)
	}
	if itemSpec.MaxListLen != 100 {
		t.Fatalf("max list len: got %d", itemSpec.MaxListLen)
	}

	spec, ok = m["ServerHTTPOrder"]
	if !ok {
		t.Fatal("expected ServerHTTPOrder spec")
	}
	orderCfg, ok := spec.(*srvhttp.Config[*orderhttp.API])
	if !ok {
		t.Fatalf("unexpected order spec type %T", spec)
	}
	if orderCfg.Label.GetValue() != "orders" {
		t.Fatalf("order label: got %q", orderCfg.Label.GetValue())
	}
}

func TestAppResources_Bootstrap(t *testing.T) {
	t.Setenv("DEMO_APP_SHUTDOWN_TIMEOUT", "5s")
	t.Setenv("DEMO_SERVER_HTTP_ITEM_LABEL", "demo")
	t.Setenv("DEMO_SERVER_HTTP_ITEM_HOST", "127.0.0.1")
	t.Setenv("DEMO_SERVER_HTTP_ITEM_PORT", "8080")
	t.Setenv("DEMO_SERVICE_ITEM_MAX_LIST_LEN", "100")
	t.Setenv("DEMO_SERVER_HTTP_ORDER_LABEL", "orders")
	t.Setenv("DEMO_SERVER_HTTP_ORDER_HOST", "127.0.0.1")
	t.Setenv("DEMO_SERVER_HTTP_ORDER_PORT", "8081")

	var ar AppResources
	reg, err := app.Bootstrap(&ar, app.Pipeline{
		Registry:   res.New(),
		EnvPrefix:  envPrefix,
		Transforms: []res.TransformFunc{},
	})
	if err != nil {
		t.Fatal(err)
	}

	n := 0
	reg.WalkEntries(func(_ res.Entry) bool {
		n++
		return true
	})
	if n == 0 {
		t.Fatal("expected resources")
	}
}

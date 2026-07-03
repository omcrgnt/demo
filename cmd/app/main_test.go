package main

import (
	"reflect"
	"testing"

	"github.com/omcrgnt/app"
	"github.com/omcrgnt/demo/internal/domain/service/item"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/res/unique"
)

func testRegistry(t *testing.T) *unique.Registry {
	t.Helper()
	return unique.Global()
}

func TestAppResources_Bootstrap(t *testing.T) {
	t.Setenv("DEMO_SERVER_HTTP_ITEM_LABEL", "demo")
	t.Setenv("DEMO_SERVER_HTTP_ITEM_HOST", "127.0.0.1")
	t.Setenv("DEMO_SERVER_HTTP_ITEM_PORT", "18080")
	t.Setenv("DEMO_SERVICE_ITEM_MAX_LIST_LEN", "100")
	t.Setenv("DEMO_SERVER_HTTP_ORDER_LABEL", "orders")
	t.Setenv("DEMO_SERVER_HTTP_ORDER_HOST", "127.0.0.1")
	t.Setenv("DEMO_SERVER_HTTP_ORDER_PORT", "18081")
	t.Setenv("DEMO_OPS_HTTP_LABEL", "ops")
	t.Setenv("DEMO_OPS_HTTP_HOST", "127.0.0.1")
	t.Setenv("DEMO_OPS_HTTP_PORT", "19090")

	reg := testRegistry(t)
	var ar _appResources

	if _, err := app.Bootstrap(&ar, app.Pipeline{
		Registry:  reg,
		EnvPrefix: envPrefix,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := reg.GetOneByType(reflect.TypeOf((*item.Service)(nil))); err != nil {
		t.Fatalf("item service: %v", err)
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

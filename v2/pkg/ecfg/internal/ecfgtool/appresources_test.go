package ecfgtool

import (
	"testing"
)

func TestCollectTemplateEntries_appResources(t *testing.T) {
	entries, err := CollectTemplateEntries(
		"github.com/omcrgnt/demo/v2/cmd/demo",
		"AppResources",
		"DEMO",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 8 {
		t.Fatalf("entries: %+v", entries)
	}
	keys := map[string]bool{}
	for _, e := range entries {
		keys[e.EnvKey] = true
		if e.Usage == "" {
			t.Fatalf("empty usage for %s", e.EnvKey)
		}
	}
	for _, want := range []string{
		"DEMO_APP_SHUTDOWN_TIMEOUT",
		"DEMO_SERVER_HTTP_ITEM_LABEL",
		"DEMO_SERVER_HTTP_ITEM_HOST",
		"DEMO_SERVER_HTTP_ITEM_PORT",
		"DEMO_SERVICE_ITEM_MAX_LIST_LEN",
		"DEMO_SERVER_HTTP_ORDER_LABEL",
		"DEMO_SERVER_HTTP_ORDER_HOST",
		"DEMO_SERVER_HTTP_ORDER_PORT",
	} {
		if !keys[want] {
			t.Fatalf("missing %s in %+v", want, entries)
		}
	}
}

package config_test

import (
	"reflect"
	"testing"

	"github.com/omcrgnt/builder"
	"github.com/omcrgnt/demo/internal/config"
	"github.com/omcrgnt/ecfg"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/sdi"
)

func TestAppConfig_Parse(t *testing.T) {
	t.Setenv("DEMO_HTTP_SERVER_LABEL", "demo")
	t.Setenv("DEMO_HTTP_SERVER_HOST", "127.0.0.1")
	t.Setenv("DEMO_HTTP_SERVER_PORT", "8080")

	cfg, err := ecfg.Parse[config.AppConfig](ecfg.WithPrefix(config.Prefix))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPServer == nil {
		t.Fatal("expected HTTPServer to be populated")
	}
	if cfg.HTTPServer.Label.GetValue() != "demo" {
		t.Fatalf("label: got %q", cfg.HTTPServer.Label.GetValue())
	}
	if cfg.HTTPServer.Host.GetValue() != "127.0.0.1" {
		t.Fatalf("host: got %q", cfg.HTTPServer.Host.GetValue())
	}
	if cfg.HTTPServer.Port.GetValue() != 8080 {
		t.Fatalf("port: got %d", cfg.HTTPServer.Port.GetValue())
	}
}

func TestAppConfig_Resolve(t *testing.T) {
	t.Setenv("DEMO_HTTP_SERVER_LABEL", "demo")
	t.Setenv("DEMO_HTTP_SERVER_HOST", "127.0.0.1")
	t.Setenv("DEMO_HTTP_SERVER_PORT", "8080")

	cfg, err := ecfg.Parse[config.AppConfig](ecfg.WithPrefix(config.Prefix))
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.Build(cfg, res.Default); err != nil {
		t.Fatal(err)
	}
	if err := sdi.Resolve(res.Default); err != nil {
		t.Fatal(err)
	}

	n := 0
	res.Walk(func(_ reflect.Type, _ any) bool {
		n++
		return true
	})
	if n == 0 {
		t.Fatal("expected resources")
	}
}

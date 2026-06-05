package config_test

import (
	"testing"

	"github.com/omcrgnt/demo/internal/config"
	"github.com/omcrgnt/demo/internal/wiring"
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

	if err := res.Build(cfg.ResSource()); err != nil {
		t.Fatal(err)
	}

	source, err := wiring.FromRegistry()
	if err != nil {
		t.Fatal(err)
	}

	di := sdi.New()
	if err := di.Resolve(source); err != nil {
		t.Fatal(err)
	}
	if len(di.Resources()) == 0 {
		t.Fatal("expected resources")
	}
}

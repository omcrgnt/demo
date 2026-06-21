package main_test

import (
	"testing"

	"github.com/omcrgnt/demo/v2/pkg/ecfg"
	"github.com/omcrgnt/demo/v2/pkg/ecfg/internal/testdata"
)

func TestExampleParse(t *testing.T) {
	t.Setenv("APP_SERVER_LABEL", "demo")
	cfg, err := ecfg.Parse[testdata.AppConfig](ecfg.WithPrefix("APP"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Label != "demo" {
		t.Fatalf("label: %q", cfg.Server.Label)
	}
}

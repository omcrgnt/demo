package ecfg_test

import (
	"errors"
	"testing"

	"github.com/omcrgnt/demo/v2/pkg/ecfg"
	"github.com/omcrgnt/demo/v2/pkg/ecfg/internal/testdata"
)

func TestPublicParse_delegates(t *testing.T) {
	t.Setenv("APP_SERVER_LABEL", "ok")
	cfg, err := ecfg.Parse[testdata.AppConfig](ecfg.WithPrefix("APP"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Label != "ok" {
		t.Fatalf("got %q", cfg.Server.Label)
	}
}

func TestPublicErrors_alias(t *testing.T) {
	t.Setenv("APP_SERVER_LABEL", "")
	_, err := ecfg.Parse[testdata.AppConfig](ecfg.WithPrefix("APP"))
	if !errors.Is(err, ecfg.ErrEmptyEnv) {
		t.Fatalf("got %v", err)
	}
}

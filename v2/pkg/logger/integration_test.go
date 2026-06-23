package logger_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/omcrgnt/demo/v2/pkg/builder"
	"github.com/omcrgnt/demo/v2/pkg/logger"
	"github.com/omcrgnt/demo/v2/pkg/res"
	"github.com/omcrgnt/demo/v2/pkg/res/restest"
	"github.com/omcrgnt/demo/v2/pkg/sdi"
	loggerv1 "github.com/omcrgnt/proto/gen/go/logger/v1"
)

func setupUseDefaults(t *testing.T) {
	t.Helper()
	restest.ResetGlobal()
	_ = res.AddToGlobalWithTags(logger.DefaultLogConfig(), res.TagReplaceable)
	_ = res.AddToGlobalWithTags(logger.DefaultStdoutConfig(), res.TagReplaceable)
	if err := builder.Build(res.Global()); err != nil {
		t.Fatal(err)
	}
}

func TestIntegration_outputOnlyOverride(t *testing.T) {
	setupUseDefaults(t)
	dir := t.TempDir()
	path := dir + "/output-only.log"

	if err := res.Global().Add(logger.OutputFileConfig{Path: path}); err != nil { //nolint:forbidigo // simulates ecfg.Register
		t.Fatal(err)
	}
	if err := builder.Build(res.Global()); err != nil {
		t.Fatal(err)
	}

	if err := sdi.Resolve(res.Global()); err != nil {
		t.Fatal(err)
	}

	logger.Info(context.Background(), "output-only")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "output-only") {
		t.Fatalf("expected log in file, got %q", data)
	}
}

func TestIntegration_dedupAndResolve(t *testing.T) {
	setupUseDefaults(t)
	dir := t.TempDir()
	path := dir + "/integration.log"

	if err := res.Global().Add(logger.OutputFileConfig{Path: path}); err != nil { //nolint:forbidigo // simulates ecfg.Register
		t.Fatal(err)
	}
	if err := res.Global().Add(logger.Config{
		Level:  &loggerv1.Level{Value: "info"},
		Format: &loggerv1.Format{Value: "json"},
	}); err != nil { //nolint:forbidigo // simulates ecfg.Register
		t.Fatal(err)
	}
	if err := builder.Build(res.Global()); err != nil {
		t.Fatal(err)
	}

	if err := sdi.Resolve(res.Global()); err != nil {
		t.Fatal(err)
	}

	logger.Info(context.Background(), "after-resolve")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "after-resolve") {
		t.Fatalf("expected log in file, got %q", data)
	}
}

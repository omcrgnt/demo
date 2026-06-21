package logger

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/omcrgnt/demo/v2/pkg/builder"
	"github.com/omcrgnt/demo/v2/pkg/res"
	"github.com/omcrgnt/demo/v2/pkg/sdi"
	loggerv1 "github.com/omcrgnt/proto/gen/go/logger/v1"
)

type testOutput struct {
	buf bytes.Buffer
}

func (t *testOutput) Write(p []byte) (int, error) { return t.buf.Write(p) }

func (t *testOutput) markerLoggerOutput() {}

// setupUseDefaults clears res and registers system Output + Log, as logger/use init does.
func setupUseDefaults() {
	res.ResetDefault()
	active = noopLogger{}
	_ = res.AddWithTags(DefaultLogConfig(), res.TagReplaceable)
	_ = res.AddWithTags(DefaultStdoutConfig(), res.TagReplaceable)
	_ = builder.Build(res.Default)
}

func TestNoop_withoutRegistry(t *testing.T) {
	res.ResetDefault()
	active = noopLogger{}

	Info(context.Background(), "silent")
	Default().Info(context.Background(), "silent")
}

func TestConfig_build(t *testing.T) {
	res.ResetDefault()
	active = noopLogger{}
	if err := res.Add(OutputStdoutConfig{}); err != nil { //nolint:forbidigo // simulates ecfg.Register
		t.Fatal(err)
	}
	if err := res.Add(Config{
		Level:  &loggerv1.Level{Value: "info"},
		Format: &loggerv1.Format{Value: "json"},
	}); err != nil { //nolint:forbidigo // simulates ecfg.Register
		t.Fatal(err)
	}
	if err := builder.Build(res.Default); err != nil {
		t.Fatal(err)
	}
	if err := sdi.Resolve(res.Default); err != nil {
		t.Fatal(err)
	}

	Info(context.Background(), "hello")
}

func TestDefault_fromRegistry(t *testing.T) {
	out := &testOutput{}
	res.ResetDefault()
	active = noopLogger{}
	_ = res.AddWithTags(out, res.TagReplaceable)
	_ = res.AddWithTags(DefaultLogConfig(), res.TagReplaceable)
	if err := builder.Build(res.Default); err != nil {
		t.Fatal(err)
	}
	if err := sdi.Resolve(res.Default); err != nil {
		t.Fatal(err)
	}

	Default().Info(context.Background(), "instance")
	if !strings.Contains(out.buf.String(), "instance") {
		t.Fatalf("expected instance log, got %q", out.buf.String())
	}
}

func TestLogImpl_Inject(t *testing.T) {
	out := &testOutput{}
	raw, err := Config{
		Level:  &loggerv1.Level{Value: "debug"},
		Format: &loggerv1.Format{Value: "text"},
	}.Build()
	if err != nil {
		t.Fatal(err)
	}
	l := raw.(*logger)
	l.Inject([]any{out})

	l.Debug(context.Background(), "debug-msg")
	if !strings.Contains(out.buf.String(), "debug-msg") {
		t.Fatalf("expected debug in output, got %q", out.buf.String())
	}
}

func TestOutputStdoutConfig(t *testing.T) {
	res.ResetDefault()
	active = noopLogger{}

	if err := res.Add(OutputStdoutConfig{}); err != nil { //nolint:forbidigo // simulates ecfg.Register
		t.Fatal(err)
	}
	if err := res.Add(Config{
		Level:  &loggerv1.Level{Value: "info"},
		Format: &loggerv1.Format{Value: "text"},
	}); err != nil { //nolint:forbidigo // simulates ecfg.Register
		t.Fatal(err)
	}
	if err := builder.Build(res.Default); err != nil {
		t.Fatal(err)
	}
	if err := sdi.Resolve(res.Default); err != nil {
		t.Fatal(err)
	}

	Info(context.Background(), "stdout-line")
}

func TestOutputFileConfig(t *testing.T) {
	res.ResetDefault()
	active = noopLogger{}

	dir := t.TempDir()
	path := dir + "/app.log"

	if err := res.Add(OutputFileConfig{Path: path}); err != nil { //nolint:forbidigo // simulates ecfg.Register
		t.Fatal(err)
	}
	if err := res.Add(Config{
		Level:  &loggerv1.Level{Value: "info"},
		Format: &loggerv1.Format{Value: "text"},
	}); err != nil { //nolint:forbidigo // simulates ecfg.Register
		t.Fatal(err)
	}
	if err := builder.Build(res.Default); err != nil {
		t.Fatal(err)
	}
	if err := sdi.Resolve(res.Default); err != nil {
		t.Fatal(err)
	}

	Info(context.Background(), "file-line")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "file-line") {
		t.Fatalf("file content: %q", data)
	}
}

func TestOutputFileConfig_emptyPath(t *testing.T) {
	_, err := OutputFileConfig{}.Build()
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestInfo_afterResolve(t *testing.T) {
	setupUseDefaults()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = old })

	if err := sdi.Resolve(res.Default); err != nil {
		t.Fatal(err)
	}

	Info(context.Background(), "bootstrap-test")

	done := make(chan []byte, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.Bytes()
	}()
	_ = w.Close()
	out := <-done
	if !strings.Contains(string(out), "bootstrap-test") {
		t.Fatalf("stdout: %q", out)
	}
}

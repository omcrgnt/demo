package gen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/omcrgnt/demo/v2/pkg/ecfg/internal/gen"
)

func TestRun_appResources(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env.template")
	if err := gen.Run("AppResources", "github.com/omcrgnt/demo/v2/cmd/demo", "DEMO", gen.Options{
		TemplatePath: out,
		MarkdownPath: filepath.Join(dir, "env.md"),
	}); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "DEMO_SERVER_HTTP_ITEM_LABEL=") {
		t.Fatalf("body: %s", body)
	}
}

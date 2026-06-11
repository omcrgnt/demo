package item_test

import (
	"testing"

	"github.com/omcrgnt/builder"
	"github.com/omcrgnt/demo/internal/data/sync/item/memory"
	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/domain/service/item"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/sdi"
)

func resolveService(t *testing.T) *item.Service {
	t.Helper()

	source := struct {
		ItemRepo memory.Config
		Service  item.Config
	}{
		ItemRepo: memory.Config{},
		Service:  item.Config{},
	}

	if err := builder.Build(source, res.Default); err != nil {
		t.Fatal(err)
	}
	if err := sdi.Resolve(res.Default); err != nil {
		t.Fatal(err)
	}

	svcs := res.Find[*item.Service]()
	if len(svcs) != 1 {
		t.Fatal("service not found in resources")
	}
	return svcs[0]
}

func TestService(t *testing.T) {
	svc := resolveService(t)

	t.Run("CRUD", func(t *testing.T) {
		created, err := svc.Create("alpha")
		if err != nil {
			t.Fatal(err)
		}

		got, err := svc.Get(created.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.Title != "alpha" {
			t.Fatalf("unexpected title: %q", got.Title)
		}

		updated, err := svc.Update(created.ID, "beta")
		if err != nil {
			t.Fatal(err)
		}
		if updated.Title != "beta" {
			t.Fatalf("unexpected title: %q", updated.Title)
		}

		items, err := svc.List()
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}

		if err := svc.Delete(created.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := svc.Get(created.ID); err != domain.ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("InvalidInput", func(t *testing.T) {
		if _, err := svc.Create("  "); err != domain.ErrInvalidInput {
			t.Fatalf("expected ErrInvalidInput, got %v", err)
		}
		if _, err := svc.Get(""); err != domain.ErrInvalidInput {
			t.Fatalf("expected ErrInvalidInput, got %v", err)
		}
	})
}

package service_test

import (
	"testing"

	"github.com/omcrgnt/builder"
	"github.com/omcrgnt/demo/internal/service"
	"github.com/omcrgnt/demo/internal/store"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/sdi"
)

func resolveService(t *testing.T) *service.Service {
	t.Helper()

	source := struct {
		Store   store.Config
		Service service.Config
	}{
		Store:   store.Config{},
		Service: service.Config{},
	}

	if err := builder.Build(source, res.Default); err != nil {
		t.Fatal(err)
	}
	if err := sdi.Resolve(res.Default); err != nil {
		t.Fatal(err)
	}

	svcs := res.Find[*service.Service]()
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
		if _, err := svc.Get(created.ID); err != service.ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("InvalidInput", func(t *testing.T) {
		if _, err := svc.Create("  "); err != service.ErrInvalidInput {
			t.Fatalf("expected ErrInvalidInput, got %v", err)
		}
		if _, err := svc.Get(""); err != service.ErrInvalidInput {
			t.Fatalf("expected ErrInvalidInput, got %v", err)
		}
	})
}

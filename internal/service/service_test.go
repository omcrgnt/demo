package service_test

import (
	"testing"

	"github.com/omcrgnt/demo/internal/service"
	"github.com/omcrgnt/demo/internal/store"
	"github.com/omcrgnt/demo/internal/wiring"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/sdi"
)

func resolveService(t *testing.T) *service.Service {
	t.Helper()

	if err := res.Build(struct {
		Store   store.Config
		Service service.Config
	}{
		Store:   store.Config{},
		Service: service.Config{},
	}); err != nil {
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

	for _, res := range di.Resources() {
		if svc, ok := res.(*service.Service); ok {
			return svc
		}
	}
	t.Fatal("service not found in resources")
	return nil
}

func TestService_CRUD(t *testing.T) {
	svc := resolveService(t)

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
}

func TestService_InvalidInput(t *testing.T) {
	svc := resolveService(t)

	if _, err := svc.Create("  "); err != service.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if _, err := svc.Get(""); err != service.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

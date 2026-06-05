package store_test

import (
	"testing"

	"github.com/omcrgnt/demo/internal/store"
)

func TestMemoryStore_CRUD(t *testing.T) {
	s := store.NewMemoryStore()

	items, err := s.AllItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty list, got %d items", len(items))
	}

	created, err := s.AddItem("first")
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Title != "first" {
		t.Fatalf("unexpected create result: %+v", created)
	}

	got, err := s.ItemByID(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got != created {
		t.Fatalf("get mismatch: %+v vs %+v", got, created)
	}

	updated, err := s.SetItem(created.ID, "updated")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "updated" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}

	items, err = s.AllItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if err := s.RemoveItem(created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ItemByID(created.ID); err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStore_NotFound(t *testing.T) {
	s := store.NewMemoryStore()

	if _, err := s.ItemByID("missing"); err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if _, err := s.SetItem("missing", "x"); err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if err := s.RemoveItem("missing"); err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestConfig_Build(t *testing.T) {
	res, err := store.Config{}.Build()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.(*store.MemoryStore); !ok {
		t.Fatalf("expected *MemoryStore, got %T", res)
	}
}

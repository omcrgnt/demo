package memory_test

import (
	"testing"

	"github.com/omcrgnt/demo/v2/internal/data/sync/item/memory"
	"github.com/omcrgnt/demo/v2/internal/domain"
)

func TestRepo_CRUD(t *testing.T) {
	r := memory.NewRepo()

	items, err := r.AllItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty list, got %d items", len(items))
	}

	created, err := r.AddItem("first")
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Title != "first" {
		t.Fatalf("unexpected create result: %+v", created)
	}

	got, err := r.ItemByID(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got != created {
		t.Fatalf("get mismatch: %+v vs %+v", got, created)
	}

	updated, err := r.SetItem(created.ID, "updated")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "updated" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}

	items, err = r.AllItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if err := r.RemoveItem(created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := r.ItemByID(created.ID); err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRepo_NotFound(t *testing.T) {
	r := memory.NewRepo()

	if _, err := r.ItemByID("missing"); err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if _, err := r.SetItem("missing", "x"); err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if err := r.RemoveItem("missing"); err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRepoRoot_NewResource(t *testing.T) {
	res, err := memory.RepoRoot{}.NewResource()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.(*memory.Repo); !ok {
		t.Fatalf("expected *memory.Repo, got %T", res)
	}
}

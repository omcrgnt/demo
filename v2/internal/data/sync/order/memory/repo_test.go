package memory_test

import (
	"testing"

	"github.com/omcrgnt/demo/v2/internal/data/sync/order/memory"
	"github.com/omcrgnt/demo/v2/internal/domain"
)

func TestRepo_CRUD(t *testing.T) {
	r := memory.NewRepo()

	orders, err := r.AllOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(orders) != 0 {
		t.Fatalf("expected empty list, got %d orders", len(orders))
	}

	created, err := r.AddOrder("first")
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Title != "first" {
		t.Fatalf("unexpected create result: %+v", created)
	}

	got, err := r.OrderByID(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got != created {
		t.Fatalf("get mismatch: %+v vs %+v", got, created)
	}

	updated, err := r.SetOrder(created.ID, "updated")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "updated" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}

	orders, err = r.AllOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(orders) != 1 {
		t.Fatalf("expected 1 order, got %d", len(orders))
	}

	if err := r.RemoveOrder(created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := r.OrderByID(created.ID); err != domain.ErrNotFound {
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

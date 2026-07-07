package memory_test

import (
	"testing"

	"github.com/omcrgnt/demo/internal/data/sync/product/memory"
	"github.com/omcrgnt/demo/internal/domain"
)

func TestRepo_CRUD(t *testing.T) {
	r := memory.NewRepo()

	products, err := r.AllProducts()
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 0 {
		t.Fatalf("expected empty list, got %d products", len(products))
	}

	created, err := r.AddProduct("first")
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Title != "first" {
		t.Fatalf("unexpected create result: %+v", created)
	}

	got, err := r.ProductByID(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got != created {
		t.Fatalf("get mismatch: %+v vs %+v", got, created)
	}

	updated, err := r.SetProduct(created.ID, "updated")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "updated" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}

	products, err = r.AllProducts()
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}

	if err := r.RemoveProduct(created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := r.ProductByID(created.ID); err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

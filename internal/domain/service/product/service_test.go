package product_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/omcrgnt/demo/internal/data/sync/product/memory"
	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/domain/service/product"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/res/restest"
	"github.com/omcrgnt/sdi"
)

var testCtx = context.Background()

func resolveService(t *testing.T) *product.Service {
	t.Helper()

	restest.ResetGlobal()
	repo, err := memory.NewRepo().NewResource()
	if err != nil {
		t.Fatal(err)
	}
	svc, err := (&product.Service{}).NewResource()
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Global().Add(repo)
	_ = res.Global().Add(svc)

	if err := sdi.Resolve(res.Global()); err != nil {
		t.Fatal(err)
	}

	svcAny, err := res.Global().GetOneByType(reflect.TypeOf((*product.Service)(nil)))
	if err != nil {
		t.Fatal(err)
	}
	return svcAny.(*product.Service)
}

func TestService_CRUD(t *testing.T) {
	svc := resolveService(t)

	created, err := svc.Create(testCtx, "alpha")
	if err != nil {
		t.Fatal(err)
	}

	got, err := svc.Get(testCtx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "alpha" {
		t.Fatalf("unexpected title: %q", got.Title)
	}

	if _, err := svc.Create(testCtx, "  "); err != domain.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

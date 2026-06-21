package order_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/omcrgnt/demo/v2/internal/data/sync/order/memory"
	"github.com/omcrgnt/demo/v2/internal/domain"
	"github.com/omcrgnt/demo/v2/internal/domain/service/order"
	"github.com/omcrgnt/demo/v2/pkg/builder"
	"github.com/omcrgnt/demo/v2/pkg/res"
	"github.com/omcrgnt/demo/v2/pkg/sdi"
)

var testCtx = context.Background()

func resolveService(t *testing.T) *order.Service {
	t.Helper()

	res.ResetDefault()
	repo, err := memory.RepoRoot{}.NewResource()
	if err != nil {
		t.Fatal(err)
	}
	svc, err := (&order.Service{}).NewResource()
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Add(repo)
	_ = res.Add(svc)

	if err := builder.Build(res.Default); err != nil {
		t.Fatal(err)
	}
	if err := sdi.Resolve(res.Default); err != nil {
		t.Fatal(err)
	}

	svcAny, err := res.GetOneByType(reflect.TypeOf((*order.Service)(nil)))
	if err != nil {
		t.Fatal(err)
	}
	return svcAny.(*order.Service)
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

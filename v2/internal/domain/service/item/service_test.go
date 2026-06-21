package item_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/omcrgnt/demo/v2/internal/data/sync/item/memory"
	"github.com/omcrgnt/demo/v2/internal/domain"
	"github.com/omcrgnt/demo/v2/internal/domain/service/item"
	"github.com/omcrgnt/demo/v2/pkg/builder"
	"github.com/omcrgnt/demo/v2/pkg/res"
	"github.com/omcrgnt/demo/v2/pkg/sdi"
)

var testCtx = context.Background()

func resolveService(t *testing.T, spec item.Spec) *item.Service {
	t.Helper()

	res.ResetDefault()
	repo, err := memory.RepoRoot{}.NewResource()
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Add(repo)
	_ = res.Add(spec)

	if err := builder.Build(res.Default); err != nil {
		t.Fatal(err)
	}
	if err := sdi.Resolve(res.Default); err != nil {
		t.Fatal(err)
	}

	svcAny, err := res.GetOneByType(reflect.TypeOf((*item.Service)(nil)))
	if err != nil {
		t.Fatal(err)
	}
	return svcAny.(*item.Service)
}

func TestService(t *testing.T) {
	svc := resolveService(t, item.Spec{})

	t.Run("CRUD", func(t *testing.T) {
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

		updated, err := svc.Update(testCtx, created.ID, "beta")
		if err != nil {
			t.Fatal(err)
		}
		if updated.Title != "beta" {
			t.Fatalf("unexpected title: %q", updated.Title)
		}

		items, err := svc.List(testCtx)
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}

		if err := svc.Delete(testCtx, created.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := svc.Get(testCtx, created.ID); err != domain.ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("InvalidInput", func(t *testing.T) {
		if _, err := svc.Create(testCtx, "  "); err != domain.ErrInvalidInput {
			t.Fatalf("expected ErrInvalidInput, got %v", err)
		}
		if _, err := svc.Get(testCtx, ""); err != domain.ErrInvalidInput {
			t.Fatalf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestService_ListMaxLen(t *testing.T) {
	svc := resolveService(t, item.Spec{MaxListLen: 2})

	for _, title := range []string{"a", "b", "c"} {
		if _, err := svc.Create(testCtx, title); err != nil {
			t.Fatal(err)
		}
	}

	items, err := svc.List(testCtx)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

func TestMaxListLen_Validate(t *testing.T) {
	if err := item.MaxListLen(-1).Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

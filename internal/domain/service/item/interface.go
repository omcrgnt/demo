package item

import (
	"context"

	"github.com/omcrgnt/demo/internal/domain/model"
)

type ItemRepository interface {
	AllItems() ([]model.Item, error)
	ItemByID(id string) (model.Item, error)
	AddItem(title string) (model.Item, error)
	SetItem(id, title string) (model.Item, error)
	RemoveItem(id string) error
}

// ItemService is the inbound port for item use cases (HTTP and other adapters).
type ItemService interface {
	List(ctx context.Context) ([]model.Item, error)
	Get(ctx context.Context, id string) (model.Item, error)
	Create(ctx context.Context, title string) (model.Item, error)
	Update(ctx context.Context, id, title string) (model.Item, error)
	Delete(ctx context.Context, id string) error
}

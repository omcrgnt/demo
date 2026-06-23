package order

import (
	"context"

	"github.com/omcrgnt/demo/internal/domain/model"
)

type OrderRepository interface {
	AllOrders() ([]model.Order, error)
	OrderByID(id string) (model.Order, error)
	AddOrder(title string) (model.Order, error)
	SetOrder(id, title string) (model.Order, error)
	RemoveOrder(id string) error
}

// OrderService is the inbound port for order use cases.
type OrderService interface {
	List(ctx context.Context) ([]model.Order, error)
	Get(ctx context.Context, id string) (model.Order, error)
	Create(ctx context.Context, title string) (model.Order, error)
	Update(ctx context.Context, id, title string) (model.Order, error)
	Delete(ctx context.Context, id string) error
}

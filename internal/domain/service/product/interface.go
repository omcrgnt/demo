package product

import (
	"context"

	"github.com/omcrgnt/demo/internal/domain/model"
)

type ProductRepository interface {
	AllProducts() ([]model.Product, error)
	ProductByID(id string) (model.Product, error)
	AddProduct(title string) (model.Product, error)
	SetProduct(id, title string) (model.Product, error)
	RemoveProduct(id string) error
}

// ProductService is the inbound port for product use cases.
type ProductService interface {
	List(ctx context.Context) ([]model.Product, error)
	Get(ctx context.Context, id string) (model.Product, error)
	Create(ctx context.Context, title string) (model.Product, error)
	Update(ctx context.Context, id, title string) (model.Product, error)
	Delete(ctx context.Context, id string) error
}

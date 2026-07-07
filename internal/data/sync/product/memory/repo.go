package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/domain/model"
)

type Repo struct {
	mu       sync.RWMutex
	products map[string]model.Product
}

func NewRepo() *Repo {
	return &Repo{
		products: make(map[string]model.Product),
	}
}

func (*Repo) NewResource() (any, error) {
	return NewRepo(), nil
}

func (r *Repo) AllProducts() ([]model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]model.Product, 0, len(r.products))
	for _, product := range r.products {
		out = append(out, product)
	}
	return out, nil
}

func (r *Repo) ProductByID(id string) (model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[id]
	if !ok {
		return model.Product{}, domain.ErrNotFound
	}
	return product, nil
}

func (r *Repo) AddProduct(title string) (model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product := model.Product{
		ID:    uuid.NewString(),
		Title: title,
	}
	r.products[product.ID] = product
	return product, nil
}

func (r *Repo) SetProduct(id, title string) (model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, ok := r.products[id]
	if !ok {
		return model.Product{}, domain.ErrNotFound
	}
	product.Title = title
	r.products[id] = product
	return product, nil
}

func (r *Repo) RemoveProduct(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.products, id)
	return nil
}

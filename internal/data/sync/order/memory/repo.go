package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/domain/model"
)

type Repo struct {
	mu     sync.RWMutex
	orders map[string]model.Order
}

func NewRepo() *Repo {
	return &Repo{
		orders: make(map[string]model.Order),
	}
}

func (*Repo) NewResource() (any, error) {
	return NewRepo(), nil
}

func (r *Repo) AllOrders() ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]model.Order, 0, len(r.orders))
	for _, order := range r.orders {
		out = append(out, order)
	}
	return out, nil
}

func (r *Repo) OrderByID(id string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return model.Order{}, domain.ErrNotFound
	}
	return order, nil
}

func (r *Repo) AddOrder(title string) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order := model.Order{
		ID:    uuid.NewString(),
		Title: title,
	}
	r.orders[order.ID] = order
	return order, nil
}

func (r *Repo) SetOrder(id, title string) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[id]
	if !ok {
		return model.Order{}, domain.ErrNotFound
	}
	order.Title = title
	r.orders[id] = order
	return order, nil
}

func (r *Repo) RemoveOrder(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.orders, id)
	return nil
}

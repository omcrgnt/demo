package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/domain/model"
)

// RepoRoot is the AppResources wire type for the item repository.
type RepoRoot struct{}

func (RepoRoot) NewResource() (any, error) {
	return NewRepo(), nil
}

type Repo struct {
	mu    sync.RWMutex
	items map[string]model.Item
}

func NewRepo() *Repo {
	return &Repo{
		items: make(map[string]model.Item),
	}
}

func (r *Repo) AllItems() ([]model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]model.Item, 0, len(r.items))
	for _, item := range r.items {
		out = append(out, item)
	}
	return out, nil
}

func (r *Repo) ItemByID(id string) (model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return model.Item{}, domain.ErrNotFound
	}
	return item, nil
}

func (r *Repo) AddItem(title string) (model.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := model.Item{
		ID:    uuid.NewString(),
		Title: title,
	}
	r.items[item.ID] = item
	return item, nil
}

func (r *Repo) SetItem(id, title string) (model.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.items[id]
	if !ok {
		return model.Item{}, domain.ErrNotFound
	}
	item.Title = title
	r.items[id] = item
	return item, nil
}

func (r *Repo) RemoveItem(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.items, id)
	return nil
}

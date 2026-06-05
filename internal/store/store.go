package store

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/omcrgnt/demo/internal/domain"
)

var ErrNotFound = errors.New("item not found")

type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]domain.Item
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]domain.Item),
	}
}

func (s *MemoryStore) AllItems() ([]domain.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]domain.Item, 0, len(s.items))
	for _, item := range s.items {
		out = append(out, item)
	}
	return out, nil
}

func (s *MemoryStore) ItemByID(id string) (domain.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return domain.Item{}, ErrNotFound
	}
	return item, nil
}

func (s *MemoryStore) AddItem(title string) (domain.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := domain.Item{
		ID:    uuid.NewString(),
		Title: title,
	}
	s.items[item.ID] = item
	return item, nil
}

func (s *MemoryStore) SetItem(id, title string) (domain.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return domain.Item{}, ErrNotFound
	}
	item.Title = title
	s.items[id] = item
	return item, nil
}

func (s *MemoryStore) RemoveItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}

type Config struct{}

func (Config) Build() (any, error) {
	return NewMemoryStore(), nil
}

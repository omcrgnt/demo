package service

import (
	"errors"
	"strings"

	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/store"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("item not found")
)

type Repository interface {
	AllItems() ([]domain.Item, error)
	ItemByID(id string) (domain.Item, error)
	AddItem(title string) (domain.Item, error)
	SetItem(id, title string) (domain.Item, error)
	RemoveItem(id string) error
}

type Service struct {
	repo Repository
}

type Config struct{}

func (Config) Build() (any, error) {
	return &Service{}, nil
}

func (s *Service) Deps() []any {
	return []any{
		(*Repository)(nil),
	}
}

func (s *Service) Inject(args []any) {
	for _, arg := range args {
		if repo, ok := arg.(Repository); ok {
			s.repo = repo
		}
	}
}

func (s *Service) List() ([]domain.Item, error) {
	return s.repo.AllItems()
}

func (s *Service) Get(id string) (domain.Item, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Item{}, ErrInvalidInput
	}
	item, err := s.repo.ItemByID(id)
	if errors.Is(err, store.ErrNotFound) {
		return domain.Item{}, ErrNotFound
	}
	return item, err
}

func (s *Service) Create(title string) (domain.Item, error) {
	if strings.TrimSpace(title) == "" {
		return domain.Item{}, ErrInvalidInput
	}
	return s.repo.AddItem(title)
}

func (s *Service) Update(id, title string) (domain.Item, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(title) == "" {
		return domain.Item{}, ErrInvalidInput
	}
	item, err := s.repo.SetItem(id, title)
	if errors.Is(err, store.ErrNotFound) {
		return domain.Item{}, ErrNotFound
	}
	return item, err
}

func (s *Service) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidInput
	}
	if err := s.repo.RemoveItem(id); errors.Is(err, store.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}

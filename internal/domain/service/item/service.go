package item

import (
	"errors"
	"strings"

	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/domain/model"
)

type Config struct{}

func (Config) Build() (any, error) {
	return &Service{}, nil
}

func (s *Service) Deps() []any {
	return []any{
		(*ItemRepository)(nil),
	}
}

func (s *Service) Inject(args []any) {
	for _, arg := range args {
		if repo, ok := arg.(ItemRepository); ok {
			s.repo = repo
		}
	}
}

type Service struct {
	repo ItemRepository
}

func (s *Service) List() ([]model.Item, error) {
	return s.repo.AllItems()
}

func (s *Service) Get(id string) (model.Item, error) {
	if strings.TrimSpace(id) == "" {
		return model.Item{}, domain.ErrInvalidInput
	}
	item, err := s.repo.ItemByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return model.Item{}, domain.ErrNotFound
	}
	return item, err
}

func (s *Service) Create(title string) (model.Item, error) {
	if strings.TrimSpace(title) == "" {
		return model.Item{}, domain.ErrInvalidInput
	}
	return s.repo.AddItem(title)
}

func (s *Service) Update(id, title string) (model.Item, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(title) == "" {
		return model.Item{}, domain.ErrInvalidInput
	}
	item, err := s.repo.SetItem(id, title)
	if errors.Is(err, domain.ErrNotFound) {
		return model.Item{}, domain.ErrNotFound
	}
	return item, err
}

func (s *Service) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.ErrInvalidInput
	}
	if err := s.repo.RemoveItem(id); errors.Is(err, domain.ErrNotFound) {
		return domain.ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}

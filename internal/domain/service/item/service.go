package item

//go:generate go run github.com/omcrgnt/obs/cmd/obsgen -type=Service

import (
	"context"
	"errors"
	"strings"

	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/domain/model"
	"github.com/omcrgnt/logger"
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

func (s *Service) Label() string {
	return "item-service"
}

func (s *Service) List(ctx context.Context) ([]model.Item, error) {
	items, err := s.repo.AllItems()
	if err != nil {
		logger.Error(ctx, "item list failed", "err", err)
		return nil, err
	}
	logger.Info(ctx, "item list", "count", len(items))
	return items, nil
}

func (s *Service) Get(ctx context.Context, id string) (model.Item, error) {
	if strings.TrimSpace(id) == "" {
		logger.Warn(ctx, "item get rejected", "reason", "empty id")
		return model.Item{}, domain.ErrInvalidInput
	}
	item, err := s.repo.ItemByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		logger.Warn(ctx, "item get not found", "id", id)
		return model.Item{}, domain.ErrNotFound
	}
	if err != nil {
		logger.Error(ctx, "item get failed", "id", id, "err", err)
		return model.Item{}, err
	}
	logger.Info(ctx, "item get", "id", item.ID, "title", item.Title)
	return item, nil
}

func (s *Service) Create(ctx context.Context, title string) (model.Item, error) {
	if strings.TrimSpace(title) == "" {
		logger.Warn(ctx, "item create rejected", "reason", "empty title")
		return model.Item{}, domain.ErrInvalidInput
	}
	item, err := s.repo.AddItem(title)
	if err != nil {
		logger.Error(ctx, "item create failed", "title", title, "err", err)
		return model.Item{}, err
	}
	logger.Info(ctx, "item create", "id", item.ID, "title", item.Title)
	return item, nil
}

func (s *Service) Update(ctx context.Context, id, title string) (model.Item, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(title) == "" {
		logger.Warn(ctx, "item update rejected", "id", id, "reason", "empty id or title")
		return model.Item{}, domain.ErrInvalidInput
	}
	item, err := s.repo.SetItem(id, title)
	if errors.Is(err, domain.ErrNotFound) {
		logger.Warn(ctx, "item update not found", "id", id)
		return model.Item{}, domain.ErrNotFound
	}
	if err != nil {
		logger.Error(ctx, "item update failed", "id", id, "err", err)
		return model.Item{}, err
	}
	logger.Info(ctx, "item update", "id", item.ID, "title", item.Title)
	return item, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		logger.Warn(ctx, "item delete rejected", "reason", "empty id")
		return domain.ErrInvalidInput
	}
	if err := s.repo.RemoveItem(id); errors.Is(err, domain.ErrNotFound) {
		logger.Warn(ctx, "item delete not found", "id", id)
		return domain.ErrNotFound
	} else if err != nil {
		logger.Error(ctx, "item delete failed", "id", id, "err", err)
		return err
	}
	logger.Info(ctx, "item delete", "id", id)
	return nil
}

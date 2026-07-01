package item

//go:generate go run github.com/omcrgnt/sdi/cmd/sdigen
//go:generate go run github.com/omcrgnt/obs/cmd/obsgen -type=Service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/domain/model"
	"github.com/omcrgnt/logger"
)

const defaultMaxListLen = 100

// MaxListLen limits how many items [Service.List] returns (0 = use default).
type MaxListLen int

func (MaxListLen) Usage() string {
	return "Maximum items returned by List (0 = default 100)"
}

func (l MaxListLen) Validate() error {
	if l < 0 {
		return fmt.Errorf("max list len must be >= 0")
	}
	return nil
}

type deps struct {
	repo ItemRepository
}

// Service is the item domain service resource.
type Service struct {
	deps
	maxListLen int
	metrics    *serviceMetrics
}

// Spec is the item service config; [Spec.Build] returns [*Service].
type Spec struct {
	MaxListLen MaxListLen
}

func (s Spec) Build() (any, error) {
	max := int(s.MaxListLen)
	if max == 0 {
		max = defaultMaxListLen
	}
	return &Service{maxListLen: max}, nil
}

func (s *Service) Label() string {
	return "item-service"
}

func (s *Service) List(ctx context.Context) ([]model.Item, error) {
	items, err := s.repo.AllItems()
	if err != nil {
		logger.Error(ctx, "item list failed", "err", err)
		s.recordOp("list", classifyErr(err))
		return nil, err
	}
	if s.maxListLen > 0 && len(items) > s.maxListLen {
		items = items[:s.maxListLen]
	}
	logger.Info(ctx, "item list", "count", len(items), "max", s.maxListLen)
	s.recordOp("list", "ok")
	return items, nil
}

func (s *Service) Get(ctx context.Context, id string) (model.Item, error) {
	if strings.TrimSpace(id) == "" {
		logger.Warn(ctx, "item get rejected", "reason", "empty id")
		s.recordOp("get", "invalid")
		return model.Item{}, domain.ErrInvalidInput
	}
	item, err := s.repo.ItemByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		logger.Warn(ctx, "item get not found", "id", id)
		s.recordOp("get", "not_found")
		return model.Item{}, domain.ErrNotFound
	}
	if err != nil {
		logger.Error(ctx, "item get failed", "id", id, "err", err)
		s.recordOp("get", classifyErr(err))
		return model.Item{}, err
	}
	logger.Info(ctx, "item get", "id", item.ID, "title", item.Title)
	s.recordOp("get", "ok")
	return item, nil
}

func (s *Service) Create(ctx context.Context, title string) (model.Item, error) {
	if strings.TrimSpace(title) == "" {
		logger.Warn(ctx, "item create rejected", "reason", "empty title")
		s.recordOp("create", "invalid")
		return model.Item{}, domain.ErrInvalidInput
	}
	item, err := s.repo.AddItem(title)
	if err != nil {
		logger.Error(ctx, "item create failed", "title", title, "err", err)
		s.recordOp("create", classifyErr(err))
		return model.Item{}, err
	}
	logger.Info(ctx, "item create", "id", item.ID, "title", item.Title)
	s.recordOp("create", "ok")
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

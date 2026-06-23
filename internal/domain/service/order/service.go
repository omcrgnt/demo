package order

//go:generate go run github.com/omcrgnt/sdi/cmd/sdigen
//go:generate go run github.com/omcrgnt/obs/cmd/obsgen -type=Service

import (
	"context"
	"errors"
	"strings"

	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/domain/model"
	"github.com/omcrgnt/logger"
)

type deps struct {
	repo OrderRepository
}

func (s *Service) NewResource() (any, error) {
	return &Service{}, nil
}

type Service struct {
	deps
}

func (s *Service) Label() string {
	return "order-service"
}

func (s *Service) List(ctx context.Context) ([]model.Order, error) {
	orders, err := s.repo.AllOrders()
	if err != nil {
		logger.Error(ctx, "order list failed", "err", err)
		return nil, err
	}
	logger.Info(ctx, "order list", "count", len(orders))
	return orders, nil
}

func (s *Service) Get(ctx context.Context, id string) (model.Order, error) {
	if strings.TrimSpace(id) == "" {
		logger.Warn(ctx, "order get rejected", "reason", "empty id")
		return model.Order{}, domain.ErrInvalidInput
	}
	order, err := s.repo.OrderByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		logger.Warn(ctx, "order get not found", "id", id)
		return model.Order{}, domain.ErrNotFound
	}
	if err != nil {
		logger.Error(ctx, "order get failed", "id", id, "err", err)
		return model.Order{}, err
	}
	logger.Info(ctx, "order get", "id", order.ID, "title", order.Title)
	return order, nil
}

func (s *Service) Create(ctx context.Context, title string) (model.Order, error) {
	if strings.TrimSpace(title) == "" {
		logger.Warn(ctx, "order create rejected", "reason", "empty title")
		return model.Order{}, domain.ErrInvalidInput
	}
	order, err := s.repo.AddOrder(title)
	if err != nil {
		logger.Error(ctx, "order create failed", "title", title, "err", err)
		return model.Order{}, err
	}
	logger.Info(ctx, "order create", "id", order.ID, "title", order.Title)
	return order, nil
}

func (s *Service) Update(ctx context.Context, id, title string) (model.Order, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(title) == "" {
		logger.Warn(ctx, "order update rejected", "id", id, "reason", "empty id or title")
		return model.Order{}, domain.ErrInvalidInput
	}
	order, err := s.repo.SetOrder(id, title)
	if errors.Is(err, domain.ErrNotFound) {
		logger.Warn(ctx, "order update not found", "id", id)
		return model.Order{}, domain.ErrNotFound
	}
	if err != nil {
		logger.Error(ctx, "order update failed", "id", id, "err", err)
		return model.Order{}, err
	}
	logger.Info(ctx, "order update", "id", order.ID, "title", order.Title)
	return order, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		logger.Warn(ctx, "order delete rejected", "reason", "empty id")
		return domain.ErrInvalidInput
	}
	if err := s.repo.RemoveOrder(id); errors.Is(err, domain.ErrNotFound) {
		logger.Warn(ctx, "order delete not found", "id", id)
		return domain.ErrNotFound
	} else if err != nil {
		logger.Error(ctx, "order delete failed", "id", id, "err", err)
		return err
	}
	logger.Info(ctx, "order delete", "id", id)
	return nil
}

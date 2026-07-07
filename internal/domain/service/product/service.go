package product

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
	repo ProductRepository
}

func (s *Service) NewResource() (any, error) {
	return &Service{}, nil
}

type Service struct {
	deps
}

func (s *Service) Label() string {
	return "product-service"
}

func (s *Service) List(ctx context.Context) ([]model.Product, error) {
	products, err := s.repo.AllProducts()
	if err != nil {
		logger.Error(ctx, "product list failed", "err", err)
		return nil, err
	}
	logger.Info(ctx, "product list", "count", len(products))
	return products, nil
}

func (s *Service) Get(ctx context.Context, id string) (model.Product, error) {
	if strings.TrimSpace(id) == "" {
		logger.Warn(ctx, "product get rejected", "reason", "empty id")
		return model.Product{}, domain.ErrInvalidInput
	}
	product, err := s.repo.ProductByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		logger.Warn(ctx, "product get not found", "id", id)
		return model.Product{}, domain.ErrNotFound
	}
	if err != nil {
		logger.Error(ctx, "product get failed", "id", id, "err", err)
		return model.Product{}, err
	}
	logger.Info(ctx, "product get", "id", product.ID, "title", product.Title)
	return product, nil
}

func (s *Service) Create(ctx context.Context, title string) (model.Product, error) {
	if strings.TrimSpace(title) == "" {
		logger.Warn(ctx, "product create rejected", "reason", "empty title")
		return model.Product{}, domain.ErrInvalidInput
	}
	product, err := s.repo.AddProduct(title)
	if err != nil {
		logger.Error(ctx, "product create failed", "title", title, "err", err)
		return model.Product{}, err
	}
	logger.Info(ctx, "product create", "id", product.ID, "title", product.Title)
	return product, nil
}

func (s *Service) Update(ctx context.Context, id, title string) (model.Product, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(title) == "" {
		logger.Warn(ctx, "product update rejected", "id", id, "reason", "empty id or title")
		return model.Product{}, domain.ErrInvalidInput
	}
	product, err := s.repo.SetProduct(id, title)
	if errors.Is(err, domain.ErrNotFound) {
		logger.Warn(ctx, "product update not found", "id", id)
		return model.Product{}, domain.ErrNotFound
	}
	if err != nil {
		logger.Error(ctx, "product update failed", "id", id, "err", err)
		return model.Product{}, err
	}
	logger.Info(ctx, "product update", "id", product.ID, "title", product.Title)
	return product, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		logger.Warn(ctx, "product delete rejected", "reason", "empty id")
		return domain.ErrInvalidInput
	}
	if err := s.repo.RemoveProduct(id); errors.Is(err, domain.ErrNotFound) {
		logger.Warn(ctx, "product delete not found", "id", id)
		return domain.ErrNotFound
	} else if err != nil {
		logger.Error(ctx, "product delete failed", "id", id, "err", err)
		return err
	}
	logger.Info(ctx, "product delete", "id", id)
	return nil
}

package product

import (
	"context"

	grpcapi "github.com/omcrgnt/demo/internal/api/grpc"
	demov1 "github.com/omcrgnt/demo/internal/api/grpc/gen/demo/v1"
	"github.com/omcrgnt/demo/internal/domain/model"
	productsvc "github.com/omcrgnt/demo/internal/domain/service/product"
	"google.golang.org/grpc"
)

type API struct {
	demov1.UnimplementedProductServiceServer
	svc productsvc.ProductService
}

func (a *API) NewResource() (any, error) {
	return &API{}, nil
}

func (a *API) Deps() []any {
	return []any{
		(*productsvc.ProductService)(nil),
	}
}

func (a *API) Inject(args []any) {
	for _, arg := range args {
		if svc, ok := arg.(productsvc.ProductService); ok {
			a.svc = svc
		}
	}
}

func (a *API) RegisterGRPC(s *grpc.Server) {
	demov1.RegisterProductServiceServer(s, a)
}

func (a *API) ListProducts(ctx context.Context, _ *demov1.ListProductsRequest) (*demov1.ListProductsResponse, error) {
	products, err := a.svc.List(ctx)
	if err != nil {
		return nil, grpcapi.ServiceError(err)
	}
	return &demov1.ListProductsResponse{Products: toProducts(products)}, nil
}

func (a *API) GetProduct(ctx context.Context, req *demov1.GetProductRequest) (*demov1.Product, error) {
	product, err := a.svc.Get(ctx, req.GetId())
	if err != nil {
		return nil, grpcapi.ServiceError(err)
	}
	return toProduct(product), nil
}

func (a *API) CreateProduct(ctx context.Context, req *demov1.CreateProductRequest) (*demov1.Product, error) {
	product, err := a.svc.Create(ctx, req.GetTitle())
	if err != nil {
		return nil, grpcapi.ServiceError(err)
	}
	return toProduct(product), nil
}

func (a *API) UpdateProduct(ctx context.Context, req *demov1.UpdateProductRequest) (*demov1.Product, error) {
	product, err := a.svc.Update(ctx, req.GetId(), req.GetTitle())
	if err != nil {
		return nil, grpcapi.ServiceError(err)
	}
	return toProduct(product), nil
}

func (a *API) DeleteProduct(ctx context.Context, req *demov1.DeleteProductRequest) (*demov1.DeleteProductResponse, error) {
	if err := a.svc.Delete(ctx, req.GetId()); err != nil {
		return nil, grpcapi.ServiceError(err)
	}
	return &demov1.DeleteProductResponse{}, nil
}

func toProduct(product model.Product) *demov1.Product {
	return &demov1.Product{
		Id:    product.ID,
		Title: product.Title,
	}
}

func toProducts(products []model.Product) []*demov1.Product {
	if products == nil {
		return []*demov1.Product{}
	}
	out := make([]*demov1.Product, 0, len(products))
	for _, product := range products {
		out = append(out, toProduct(product))
	}
	return out
}

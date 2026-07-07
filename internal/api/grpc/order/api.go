package order

import (
	"context"

	grpcapi "github.com/omcrgnt/demo/internal/api/grpc"
	demov1 "github.com/omcrgnt/demo/internal/api/grpc/gen/demo/v1"
	"github.com/omcrgnt/demo/internal/domain/model"
	ordersvc "github.com/omcrgnt/demo/internal/domain/service/order"
	"google.golang.org/grpc"
)

type API struct {
	demov1.UnimplementedOrderServiceServer
	svc ordersvc.OrderService
}

func (a *API) NewResource() (any, error) {
	return &API{}, nil
}

func (a *API) Deps() []any {
	return []any{
		(*ordersvc.OrderService)(nil),
	}
}

func (a *API) Inject(args []any) {
	for _, arg := range args {
		if svc, ok := arg.(ordersvc.OrderService); ok {
			a.svc = svc
		}
	}
}

func (a *API) RegisterGRPC(s *grpc.Server) {
	demov1.RegisterOrderServiceServer(s, a)
}

func (a *API) ListOrders(ctx context.Context, _ *demov1.ListOrdersRequest) (*demov1.ListOrdersResponse, error) {
	orders, err := a.svc.List(ctx)
	if err != nil {
		return nil, grpcapi.ServiceError(err)
	}
	return &demov1.ListOrdersResponse{Orders: toOrders(orders)}, nil
}

func (a *API) GetOrder(ctx context.Context, req *demov1.GetOrderRequest) (*demov1.Order, error) {
	order, err := a.svc.Get(ctx, req.GetId())
	if err != nil {
		return nil, grpcapi.ServiceError(err)
	}
	return toOrder(order), nil
}

func (a *API) CreateOrder(ctx context.Context, req *demov1.CreateOrderRequest) (*demov1.Order, error) {
	order, err := a.svc.Create(ctx, req.GetTitle())
	if err != nil {
		return nil, grpcapi.ServiceError(err)
	}
	return toOrder(order), nil
}

func (a *API) UpdateOrder(ctx context.Context, req *demov1.UpdateOrderRequest) (*demov1.Order, error) {
	order, err := a.svc.Update(ctx, req.GetId(), req.GetTitle())
	if err != nil {
		return nil, grpcapi.ServiceError(err)
	}
	return toOrder(order), nil
}

func (a *API) DeleteOrder(ctx context.Context, req *demov1.DeleteOrderRequest) (*demov1.DeleteOrderResponse, error) {
	if err := a.svc.Delete(ctx, req.GetId()); err != nil {
		return nil, grpcapi.ServiceError(err)
	}
	return &demov1.DeleteOrderResponse{}, nil
}

func toOrder(order model.Order) *demov1.Order {
	return &demov1.Order{
		Id:    order.ID,
		Title: order.Title,
	}
}

func toOrders(orders []model.Order) []*demov1.Order {
	if orders == nil {
		return []*demov1.Order{}
	}
	out := make([]*demov1.Order, 0, len(orders))
	for _, order := range orders {
		out = append(out, toOrder(order))
	}
	return out
}

package bundle

import (
	ordergrpc "github.com/omcrgnt/demo/internal/api/grpc/order"
	productgrpc "github.com/omcrgnt/demo/internal/api/grpc/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Bundle struct {
	order   *ordergrpc.API
	product *productgrpc.API
}

func (b *Bundle) NewResource() (any, error) {
	return &Bundle{}, nil
}

func (b *Bundle) Deps() []any {
	return []any{
		(*ordergrpc.API)(nil),
		(*productgrpc.API)(nil),
	}
}

func (b *Bundle) Inject(args []any) {
	for _, arg := range args {
		switch v := arg.(type) {
		case *ordergrpc.API:
			b.order = v
		case *productgrpc.API:
			b.product = v
		}
	}
}

func (b *Bundle) RegisterGRPC(s *grpc.Server) {
	b.order.RegisterGRPC(s)
	b.product.RegisterGRPC(s)
	reflection.Register(s)
}

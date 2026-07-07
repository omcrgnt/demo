//go:generate go run github.com/omcrgnt/ecfg/cmd/ecfg-gen -type _appResources -pkg github.com/omcrgnt/demo/cmd/app -prefix DEMO -template ../../.env.template -md ../../env.md

package main

import (
	"log"

	_ "github.com/omcrgnt/meta/core/use"

	"github.com/omcrgnt/app"
	grpcbundle "github.com/omcrgnt/demo/internal/api/grpc/bundle"
	ordergrpc "github.com/omcrgnt/demo/internal/api/grpc/order"
	productgrpc "github.com/omcrgnt/demo/internal/api/grpc/product"
	handleritem "github.com/omcrgnt/demo/internal/api/http/item"
	handlerorder "github.com/omcrgnt/demo/internal/api/http/order"
	repoitem "github.com/omcrgnt/demo/internal/data/sync/item/memory"
	repoorder "github.com/omcrgnt/demo/internal/data/sync/order/memory"
	repoproduct "github.com/omcrgnt/demo/internal/data/sync/product/memory"
	serviceitem "github.com/omcrgnt/demo/internal/domain/service/item"
	serviceorder "github.com/omcrgnt/demo/internal/domain/service/order"
	serviceproduct "github.com/omcrgnt/demo/internal/domain/service/product"
	"github.com/omcrgnt/obs"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/res/unique"
	srvgrpc "github.com/omcrgnt/srv-grpc"
	srvhttp "github.com/omcrgnt/srv-http"
)

const envPrefix = "DEMO"

type _appResources struct {
	ServerHTTPItem *srvhttp.Server[*handleritem.API] `ecfg:"SERVER_HTTP_ITEM"`
	APIItem        *handleritem.API
	ServiceItem    *serviceitem.Service `ecfg:"SERVICE_ITEM"`
	RepoItem       *repoitem.Repo

	ServerHTTPOrder *srvhttp.Server[*handlerorder.API] `ecfg:"SERVER_HTTP_ORDER"`
	APIOrder        *handlerorder.API
	ServiceOrder    *serviceorder.Service
	RepoOrder       *repoorder.Repo

	ServerGRPC     *srvgrpc.Server[*grpcbundle.Bundle] `ecfg:"SERVER_GRPC"`
	Bundle         *grpcbundle.Bundle
	GRPCAPIOrder   *ordergrpc.API
	GRPCAPIProduct *productgrpc.API
	ServiceProduct *serviceproduct.Service
	RepoProduct    *repoproduct.Repo
}

var appResources _appResources

func main() {
	pipeline := app.Pipeline{
		Registry:   unique.Global(),
		EnvPrefix:  envPrefix,
		Transforms: []res.TransformFunc{obs.ApplyTransform},
	}

	if err := app.Run(&appResources, pipeline); err != nil {
		log.Fatal(err)
	}
}

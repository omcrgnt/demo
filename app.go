//go:generate go run github.com/omcrgnt/ecfg/cmd/ecfg-gen -type AppResources -pkg github.com/omcrgnt/demo -prefix DEMO -template .env.template -md env.md

package main

import (
	"log"

	_ "github.com/omcrgnt/app/use"
	_ "github.com/omcrgnt/logger/use"
	_ "github.com/omcrgnt/ops/metrics/use"
	_ "github.com/omcrgnt/ops/transport/http/use"
	_ "github.com/omcrgnt/srv-http/use"
	_ "github.com/omcrgnt/telemetry/use"

	"github.com/omcrgnt/demo/internal/api/http"
	orderhttp "github.com/omcrgnt/demo/internal/api/http/order"
	"github.com/omcrgnt/demo/internal/data/sync/item/memory"
	ordermemory "github.com/omcrgnt/demo/internal/data/sync/order/memory"
	"github.com/omcrgnt/demo/internal/domain/service/order"
	"github.com/omcrgnt/app"
	"github.com/omcrgnt/obs"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/res/unique"
)

const envPrefix = "DEMO"

type AppResources struct {
	RepoItem        memory.RepoRoot
	ServiceItem     serviceItemWire `ecfg:"SERVICE_ITEM"`
	RepoOrder       ordermemory.RepoRoot
	ServiceOrder    *order.Service
	ServerHTTPOps   serverOpsHTTPWire `ecfg:"OPS_HTTP"`
	ServerHTTPItem  serverHTTPItemWire `ecfg:"SERVER_HTTP_ITEM"`
	APIItem         *http.API
	ServerHTTPOrder serverHTTPOrderWire `ecfg:"SERVER_HTTP_ORDER"`
	APIOrder        *orderhttp.API
}

var appResources AppResources

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

//go:generate go run github.com/omcrgnt/ecfg/cmd/ecfg-gen -type AppResources -pkg github.com/omcrgnt/demo/cmd/demo -prefix DEMO -template ../../.env.template -md ../../env.md

package main

import (
	"log"

	_ "github.com/omcrgnt/logger/use"
	_ "github.com/omcrgnt/telemetry/use"

	"github.com/omcrgnt/demo/internal/api/http"
	orderhttp "github.com/omcrgnt/demo/internal/api/http/order"
	"github.com/omcrgnt/demo/internal/data/sync/item/memory"
	ordermemory "github.com/omcrgnt/demo/internal/data/sync/order/memory"
	"github.com/omcrgnt/demo/internal/domain/service/item"
	"github.com/omcrgnt/demo/internal/domain/service/order"
	"github.com/omcrgnt/app"
	"github.com/omcrgnt/obs"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/runner"
	srvhttp "github.com/omcrgnt/srv-http"
)

const envPrefix = "DEMO"

type AppResources struct {
	App             *app.App `ecfg:"APP"`
	Runner          *runner.Runner
	RepoItem        memory.RepoRoot
	ServiceItem     *item.Service `ecfg:"SERVICE_ITEM"`
	RepoOrder       ordermemory.RepoRoot
	ServiceOrder    *order.Service
	Metrics         http.Metrics
	ServerHTTPItem  *http.Server `ecfg:"SERVER_HTTP_ITEM"`
	APIItem         *http.API
	ServerHTTPOrder *srvhttp.Config[*orderhttp.API] `ecfg:"SERVER_HTTP_ORDER"`
	APIOrder        *orderhttp.API
}

var appResources AppResources

func main() {
	pipeline := app.Pipeline{
		Registry:   res.Global(),
		EnvPrefix:  envPrefix,
		Transforms: []res.TransformFunc{obs.ApplyTransform},
	}

	if err := app.Run(&appResources, pipeline); err != nil {
		log.Fatal(err)
	}
}

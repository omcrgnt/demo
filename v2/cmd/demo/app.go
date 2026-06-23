//go:generate go run github.com/omcrgnt/demo/v2/pkg/ecfg/cmd/ecfg-gen -type AppResources -pkg github.com/omcrgnt/demo/v2/cmd/demo -prefix DEMO -template ../../.env.template -md ../../env.md

package main

import (
	"log"

	"github.com/omcrgnt/demo/v2/internal/api/http"
	orderhttp "github.com/omcrgnt/demo/v2/internal/api/http/order"
	"github.com/omcrgnt/demo/v2/internal/data/sync/item/memory"
	ordermemory "github.com/omcrgnt/demo/v2/internal/data/sync/order/memory"
	"github.com/omcrgnt/demo/v2/internal/domain/service/item"
	"github.com/omcrgnt/demo/v2/internal/domain/service/order"
	"github.com/omcrgnt/demo/v2/pkg/app"
	"github.com/omcrgnt/demo/v2/pkg/obs"
	"github.com/omcrgnt/demo/v2/pkg/res"
	"github.com/omcrgnt/demo/v2/pkg/runner"
	srvhttp "github.com/omcrgnt/demo/v2/pkg/srv-http"
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

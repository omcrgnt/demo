//go:generate go run github.com/omcrgnt/ecfg/cmd/ecfg-gen -type AppConfig -pkg github.com/omcrgnt/demo/internal/config -prefix DEMO -o ../../env.template

package config

import (
	srvhttp "github.com/omcrgnt/srv-http"

	"github.com/omcrgnt/demo/internal/api/http"
	"github.com/omcrgnt/demo/internal/data/sync/item/memory"
	"github.com/omcrgnt/demo/internal/domain/service/item"
)

const Prefix = "DEMO"

type AppConfig struct {
	ItemRepo   memory.Config              `ecfg:"ITEM_REPO"`
	ItemService item.Config               `ecfg:"ITEM_SERVICE"`
	HTTP       http.Config                `ecfg:"HTTP"`
	Metrics    http.MetricsConfig         `ecfg:"METRICS"`
	HTTPServer *srvhttp.Config[*http.API] `ecfg:"HTTP_SERVER"`
}

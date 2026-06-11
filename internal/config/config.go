//go:generate go run github.com/omcrgnt/ecfg/cmd/ecfg-gen -type AppConfig -pkg github.com/omcrgnt/demo/internal/config -prefix DEMO -o ../../env.template

package config

import (
	srvhttp "github.com/omcrgnt/srv-http"

	"github.com/omcrgnt/demo/internal/httpapi"
	"github.com/omcrgnt/demo/internal/service"
	"github.com/omcrgnt/demo/internal/store"
)

const Prefix = "DEMO"

type AppConfig struct {
	Store      store.Config                  `ecfg:"STORE"`
	Service    service.Config                `ecfg:"SERVICE"`
	Controller httpapi.Config                `ecfg:"CONTROLLER"`
	Metrics    httpapi.MetricsConfig         `ecfg:"METRICS"`
	HTTPServer *srvhttp.Config[*httpapi.API] `ecfg:"HTTP_SERVER"`
}

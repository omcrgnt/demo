package main

import (
	"github.com/omcrgnt/app"
	"github.com/omcrgnt/demo/internal/api/http"
	orderhttp "github.com/omcrgnt/demo/internal/api/http/order"
	"github.com/omcrgnt/demo/internal/domain/service/item"
	ophttp "github.com/omcrgnt/ops/transport/http"
	srvhttp "github.com/omcrgnt/srv-http"
)

type serviceItemWire struct{}

func (serviceItemWire) BuildConfig() (app.Materializer, error) {
	return &item.Spec{}, nil
}

type serverOpsHTTPWire struct{}

func (serverOpsHTTPWire) BuildConfig() (app.Materializer, error) {
	return &ophttp.Config{}, nil
}

type serverHTTPItemWire struct{}

func (serverHTTPItemWire) BuildConfig() (app.Materializer, error) {
	cfg := srvhttp.Config[*http.API]{}
	return &cfg, nil
}

type serverHTTPOrderWire struct{}

func (serverHTTPOrderWire) BuildConfig() (app.Materializer, error) {
	cfg := srvhttp.Config[*orderhttp.API]{}
	return &cfg, nil
}

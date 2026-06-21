package main

import (
	"fmt"
	"log"

	"github.com/omcrgnt/demo/v2/pkg/ecfg"
	"github.com/omcrgnt/demo/v2/pkg/ecfg/internal/testdata"
)

//go:generate go run github.com/omcrgnt/demo/v2/pkg/ecfg/cmd/ecfg-gen -type AppConfig -pkg github.com/omcrgnt/demo/v2/pkg/ecfg/internal/testdata -prefix APP -template env.template

func main() {
	cfg, err := ecfg.Parse[testdata.AppConfig](ecfg.WithPrefix("APP"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.Server.Label)
}

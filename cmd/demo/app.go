package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omcrgnt/builder"
	"github.com/omcrgnt/demo/internal/config"
	"github.com/omcrgnt/ecfg"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/runner"
	"github.com/omcrgnt/sdi"
)

func main() {
	cfg, err := ecfg.Parse[config.AppConfig](ecfg.WithPrefix(config.Prefix))
	if err != nil {
		log.Fatal(err)
	}

	if err := builder.Build(cfg, res.Default); err != nil {
		log.Fatal(err)
	}
	if err := sdi.Resolve(res.Default); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	r := runner.New(res.Default)
	go func() {
		if err := r.Run(ctx); err != nil {
			slog.Error("runner failed", "err", err)
			cancel()
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := r.Stop(shutdownCtx); err != nil {
		slog.Error("shutdown failed", "err", err)
		os.Exit(1)
	}
}

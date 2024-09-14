package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/mainden/stdhttp/pkg/flagx"
	"github.com/mainden/stdhttp/pkg/logx"
	"github.com/mainden/stdhttp/pkg/pubsubx"
	"github.com/mainden/stdhttp/pkg/runx"
)

func main() {
	defer logx.Close()
	ctx := logx.WithName(pubsubx.WithCancel(context.Background()), "main")
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	go runx.AwaitDone(ctx, stop)
	config := configure()
	switch flagx.GetStoredCommand() {
	case "debug":
		debug(ctx, &config.Debug)
	case "run":
		run(ctx, &config.Run)
	case "broker":
		broker(ctx, &config.Broker)
	case "list":
		list(ctx, &config.List)
	case "kill":
		kill(ctx, &config.Kill)
	default:
		panic("unknown command")
	}
}

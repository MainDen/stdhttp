package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/mainden/stdhttp/internal/configs"
	"github.com/mainden/stdhttp/internal/controllers"
	"github.com/mainden/stdhttp/internal/handlers"
	"github.com/mainden/stdhttp/pkg/httpx"
	"github.com/mainden/stdhttp/pkg/iox"
	"github.com/mainden/stdhttp/pkg/logx"
	"github.com/mainden/stdhttp/pkg/runx"
)

func broker(ctx context.Context, config *configs.StdhttpBrokerConfig) {
	ctx = logx.WithName(ctx, "broker")
	stdout, err := iox.Output(config.StdoutOutput)
	if err != nil {
		logx.FatalContext(ctx, "Error creating stdout output", "error", err)
	}
	logx.InfoContext(ctx, fmt.Sprintf("Listening on %s", config.Address))
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		logx.FatalContext(ctx, "Failed to listen", "error", err)
	}
	go runx.AwaitDone(ctx, func() { listener.Close() })

	processBrocker := controllers.NewProcessesBrokerController(config.WaitTimeout)
	http.Handle("/api/v1/processes", httpx.HandleEvent(httpx.HandleLogger(http.StripPrefix("/api/v1", handlers.NewProcessesBrokerHttpHandler(processBrocker, stdout)))))
	http.Handle("/api/v1/processes/", httpx.HandleEvent(httpx.HandleLogger(http.StripPrefix("/api/v1", handlers.NewProcessesBrokerHttpHandler(processBrocker, stdout)))))
	err = http.Serve(listener, http.DefaultServeMux)
	if err != nil && !errors.Is(err, net.ErrClosed) {
		logx.FatalContext(ctx, "Failed to serve", "error", err)
	} else {
		logx.InfoContext(ctx, "Server stopped")
	}
}

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/mainden/stdhttp/internal/configs"
	"github.com/mainden/stdhttp/internal/handlers"
	"github.com/mainden/stdhttp/pkg/httpx"
	"github.com/mainden/stdhttp/pkg/iox"
	"github.com/mainden/stdhttp/pkg/logx"
	"github.com/mainden/stdhttp/pkg/runx"
)

func debug(ctx context.Context, config *configs.StdhttpDebugConfig) {
	ctx = logx.WithName(ctx, "debug")
	stdout, err := iox.Output(config.StdoutOutput)
	if err != nil {
		logx.FatalContext(ctx, "Error creating stdout output", "error", err)
	}
	stderr, err := iox.Output(config.StderrOutput)
	if err != nil {
		logx.FatalContext(ctx, "Error creating stderr output", "error", err)
	}
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		logx.FatalContext(ctx, "Failed to listen", "error", err)
	} else {
		logx.InfoContext(ctx, fmt.Sprintf("Listening on %s", config.Address))
	}
	go runx.AwaitDone(ctx, func() { listener.Close() })

	http.Handle("/", httpx.HandleEvent(httpx.HandleLogger(handlers.NewPostTextDebugHttpHandler(stdout, stderr))))
	err = http.Serve(listener, http.DefaultServeMux)
	if err != nil && !errors.Is(err, net.ErrClosed) {
		logx.FatalContext(ctx, "Failed to serve", "error", err)
	} else {
		logx.InfoContext(ctx, "Server stopped")
	}
}

package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mainden/stdhttp/internal/clients"
	"github.com/mainden/stdhttp/internal/configs"
	"github.com/mainden/stdhttp/internal/models"
	"github.com/mainden/stdhttp/pkg/iox"
	"github.com/mainden/stdhttp/pkg/logx"
)

func kill(ctx context.Context, config *configs.StdhttpKillConfig) {
	output, err := iox.Output(config.StdoutOutput)
	if err != nil {
		logx.FatalContext(ctx, "Error creating stdout output", "error", err)
	}
	pid, err := strconv.ParseInt(config.Pattern, 0, 0)
	parsed := err == nil
	if parsed {
		err = clients.NewProcessesBrokerHttpClient(config.BrokerURL, 0).Kill(ctx, int(pid))
	} else {
		err = clients.NewProcessesBrokerHttpClient(config.BrokerURL, 0).KillMany(ctx, config.Pattern)
	}
	if err != nil && !errors.Is(err, models.ErrProcessNotFound) {
		logx.FatalContext(ctx, "Error killing processes", "error", err)
	}
	if errors.Is(err, models.ErrProcessNotFound) {
		fmt.Fprintf(output, "Process not found\n")
	} else {
		if parsed {
			fmt.Fprintf(output, "Process killed\n")
		} else {
			fmt.Fprintf(output, "Processes killed\n")
		}
	}
}

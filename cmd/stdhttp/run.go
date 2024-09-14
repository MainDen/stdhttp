package main

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/mainden/stdhttp/internal/clients"
	"github.com/mainden/stdhttp/internal/configs"
	"github.com/mainden/stdhttp/internal/handlers"
	"github.com/mainden/stdhttp/internal/models"
	"github.com/mainden/stdhttp/pkg/iox"
	"github.com/mainden/stdhttp/pkg/logx"
	"github.com/mainden/stdhttp/pkg/osx/execx"
	"github.com/mainden/stdhttp/pkg/pubsubx"
	"github.com/mainden/stdhttp/pkg/runx"
	"github.com/mainden/stdhttp/pkg/textx"
)

var (
	TopicStdoutLine    = "stdout.line"
	TopicStderrLine    = "stderr.line"
	TopicBrokerCommand = "broker.command"
)

func run(ctx context.Context, config *configs.StdhttpRunConfig) {
	if config.StdoutURL != "" {
		pubsubx.Subscribe(ctx, TopicStdoutLine, handlers.NewPostTextPubsubHandler(config.StdoutURL, "stdout"))
	}
	if config.StderrURL != "" {
		pubsubx.Subscribe(ctx, TopicStderrLine, handlers.NewPostTextPubsubHandler(config.StderrURL, "stderr"))
	}
	if config.BrokerURL != "" {
		processesClient := clients.NewProcessesBrokerHttpClient(config.BrokerURL, config.BrokerWaitTimeout)
		process := models.ProcessModel{
			Pid:         os.Getpid(),
			ClientName:  config.BrokerClientName,
			CommandName: config.CommandName,
			CommandArgs: config.CommandArgs,
			Persistent:  config.Persistent,
		}
		defer runx.Await(runx.Async(func() { processesClient.CommandLoop(ctx, process, TopicBrokerCommand) }))
		defer pubsubx.Cancel(ctx)
	}

	switch {
	case config.CommandName != "":
		runCommand(ctx, config)
	default:
		runPipe(ctx, config)
	}
}

func runCommand(ctx context.Context, config *configs.StdhttpRunConfig) {
	ctx = logx.WithName(ctx, "run")
	stdout, err := iox.Output(config.StdoutOutput)
	if err != nil {
		logx.FatalContext(ctx, "Error creating stdout output", "error", err)
	}
	defer iox.Close(stdout)
	stderr, err := iox.Output(config.StderrOutput)
	if err != nil {
		logx.FatalContext(ctx, "Error creating stderr output", "error", err)
	}
	defer iox.Close(stderr)

	if config.StdoutURL != "" {
		r, w := io.Pipe()
		go runx.AwaitDone(ctx, func() { iox.Close(w) })
		stdout = io.MultiWriter(stdout, w)
		defer runx.Await(runx.Async(func() { textx.NewTextScanner(r, TopicStdoutLine).Run(context.Background()) }))
	}
	if config.StderrURL != "" {
		r, w := io.Pipe()
		go runx.AwaitDone(ctx, func() { iox.Close(w) })
		stderr = io.MultiWriter(stderr, w)
		defer runx.Await(runx.Async(func() { textx.NewTextScanner(r, TopicStderrLine).Run(context.Background()) }))
	}

	logx.InfoContext(ctx, "Starting command", "name", config.CommandName, "args", config.CommandArgs)
	group := execx.NewProcessGroup()
	defer group.Close()
	for done := false; config.Persistent || !done; done = true {
		if done {
			logx.InfoContext(ctx, "Restarting command", "name", config.CommandName, "args", config.CommandArgs)
		}

		cmd := exec.CommandContext(ctx, config.CommandName, config.CommandArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		execx.CmdHide(cmd)

		if err := cmd.Start(); err != nil && ctx.Err() == nil {
			logx.ErrorContext(ctx, "Failed to start command", "name", config.CommandName, "args", config.CommandArgs, "error", err)
			continue
		} else if ctx.Err() == nil {
			logx.InfoContext(ctx, "Command started", "name", config.CommandName, "args", config.CommandArgs)
		} else {
			logx.InfoContext(ctx, "Command cancelled", "name", config.CommandName, "args", config.CommandArgs)
			return
		}

		if err := group.Add(cmd); err != nil {
			logx.ErrorContext(ctx, "Failed to add command to process group", "name", config.CommandName, "args", config.CommandArgs, "error", err)
		}

		if err := cmd.Wait(); err != nil && ctx.Err() == nil {
			logx.ErrorContext(ctx, "Command failed", "name", config.CommandName, "args", config.CommandArgs, "error", err)
			continue
		} else if ctx.Err() == nil {
			logx.InfoContext(ctx, "Command succeeded", "name", config.CommandName, "args", config.CommandArgs)
			continue
		} else {
			logx.InfoContext(ctx, "Command cancelled", "name", config.CommandName, "args", config.CommandArgs)
			return
		}
	}
}

func runPipe(ctx context.Context, config *configs.StdhttpRunConfig) {
	stdout, err := iox.Output(config.StdoutOutput)
	if err != nil {
		logx.FatalContext(ctx, "Error creating stdout output", "error", err)
	}
	defer iox.Close(stdout)

	if config.StdoutURL != "" {
		r, w := io.Pipe()
		go runx.AwaitDone(ctx, func() { iox.Close(w) })
		stdout = io.MultiWriter(stdout, w)
		defer runx.Await(runx.Async(func() { textx.NewTextScanner(r, TopicStdoutLine).Run(context.Background()) }))
	}

	logx.InfoContext(ctx, "Piping stdin")
	if _, err := io.Copy(stdout, iox.NewContextReader(ctx, os.Stdin)); err != nil && ctx.Err() == nil {
		logx.ErrorContext(ctx, "Pipe failed", "error", err)
	} else if ctx.Err() == nil {
		logx.InfoContext(ctx, "Pipe succeeded")
	} else {
		logx.InfoContext(ctx, "Pipe cancelled")
	}
}

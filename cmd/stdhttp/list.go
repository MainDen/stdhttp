package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mainden/stdhttp/internal/clients"
	"github.com/mainden/stdhttp/internal/configs"
	"github.com/mainden/stdhttp/internal/models"
	"github.com/mainden/stdhttp/pkg/iox"
	"github.com/mainden/stdhttp/pkg/logx"
)

func list(ctx context.Context, config *configs.StdhttpListConfig) {
	output, err := iox.Output(config.StdoutOutput)
	if err != nil {
		logx.FatalContext(ctx, "Error creating stdout output", "error", err)
	}
	processes, err := clients.NewProcessesBrokerHttpClient(config.BrokerURL, 0).List(ctx)
	if err != nil {
		logx.FatalContext(ctx, "Error listing processes", "error", err)
	}
	sort.Slice(processes, func(i, j int) bool { return processes[i].Pid < processes[j].Pid })
	_, err = output.Write([]byte(listProcessesFormat(processes)))
	if err != nil {
		logx.FatalContext(ctx, "Error writing output", "error", err)
	}
}

func listProcessesFormat(processes models.ProcessModels) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%-12v %-12v %v\n", "PID", "CLIENT NAME", "COMMAND")
	for _, process := range processes {
		fmt.Fprint(&builder, listProcessFormat(process))
	}
	fmt.Fprintf(&builder, "Total: %v\n", len(processes))
	return builder.String()
}

func listProcessFormat(process models.ProcessModel) string {
	return fmt.Sprintf("%-12v %-12v %v\n", process.Pid, process.ClientName, listCommandFormat(process))
}

func listCommandFormat(process models.ProcessModel) string {
	if process.CommandName == "" {
		return "PIPE"
	}
	if len(process.CommandArgs) == 0 {
		return listQuote(process.CommandName)
	}
	return listQuote(process.CommandName) + " " + strings.Join(listQuoteSlice(process.CommandArgs), " ")
}

func listQuoteSlice(s []string) []string {
	s = append([]string(nil), s...)
	for i := range s {
		s[i] = listQuote(s[i])
	}
	return s
}

func listQuote(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	if strings.ContainsAny(s, " \"") {
		return "\"" + s + "\""
	}
	return s
}

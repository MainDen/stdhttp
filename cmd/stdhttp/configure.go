package main

import (
	"os"

	"github.com/mainden/stdhttp/internal/configs"
	"github.com/mainden/stdhttp/pkg/debugx/pex"
	"github.com/mainden/stdhttp/pkg/flagx"
	"github.com/mainden/stdhttp/pkg/httpx"
	"github.com/mainden/stdhttp/pkg/iox"
	"github.com/mainden/stdhttp/pkg/logx"
	"github.com/mainden/stdhttp/pkg/stringsx"
)

var appname string

var version string

var copyright string

var license string

var url string

func init() {
	iox.SetConsoleModeEnabled(!pex.IsGUI())
	httpx.Default = httpx.WrapHttpClient(httpx.Default, httpx.WithUserAgent(appname+"/"+version))
}

func configure() *configs.StdhttpConfig {
	var config configs.StdhttpConfig
	flagx.SetName(appname)
	flagx.SetPrefix("STDHTTP")
	flagx.SetDescription("The HTTP pipe for standard streams.")
	flagx.SetSort(true, true, true, true)
	flagx.SetOutput(iox.SelectWriter(pex.IsGUI(), pex.GUIOutput(pex.IconError(), appname), os.Stderr))

	flagx.AddHelp(iox.SelectWriter(pex.IsGUI(), pex.GUIOutput(pex.IconInformation(), appname), os.Stderr), false)
	flagx.AddVersion(version, copyright, license, url)
	flagx.AddOptEnvFunc("log-level", 'l', "LOG_LEVEL", "Sets log level.", logx.ParseLevel, flagx.WithDefaults("info"))
	flagx.AddOptEnvFunc("log-format", 'f', "LOG_FORMAT", "Sets log format.", logx.ParseFormat, flagx.WithDefaults("json"))
	flagx.AddOptEnvFunc("log-output", 'o', "{OUTPUT|FILE}", "Sets log output.", logx.ParseOutput, flagx.WithDefaults(stringsx.SelectString(pex.IsGUI(), "null", "stdout")))
	flagx.AddOpt("debug", 'd', "", "Enables debug log level.", flagx.Func(logx.ParseLevel), flagx.WithArgs("debug"))
	flagx.AddOpt("info", 'i', "", "Enables info log level.", flagx.Func(logx.ParseLevel), flagx.WithArgs("info"))
	flagx.AddOpt("warn", 'w', "", "Enables warn log level.", flagx.Func(logx.ParseLevel), flagx.WithArgs("warn"))
	flagx.AddOpt("error", 'e', "", "Enables error log level.", flagx.Func(logx.ParseLevel), flagx.WithArgs("error"))
	flagx.AddOpt("silent", 's', "", "Enables silent log level.", flagx.Func(logx.ParseLevel), flagx.WithArgs("silent"))
	flagx.AddParam("LOG_LEVEL", "The log level value. One of: debug, info, warn, error, silent.")
	flagx.AddParam("LOG_FORMAT", "The log format value. One of: json, text.")
	flagx.AddParam("OUTPUT", "The output destination. One of: stdout, stderr, null.")
	flagx.AddParam("FILE", "The file path. Example: \"file.txt\".")

	runCmd := flagx.AddCmd("run")
	runCmd.AddHelp(iox.SelectWriter(pex.IsGUI(), pex.GUIOutput(pex.IconInformation(), appname), os.Stderr), true)
	runCmd.SetDescription("Runs the command with standard streams piped to HTTP." + stringsx.SelectString(pex.IsGUI(), "", "\nIf no command is specified, standard input will be interpreted as the standard output of the command."))
	runCmd.SetShortUsage("Runs the command with standard streams piped to HTTP.")
	runCmd.AddOptEnvString("stdout-url", 'O', "URL", "Sets the URL to pipe standard output to.", &config.Run.StdoutURL)
	runCmd.AddOptEnvString("stderr-url", 'E', "URL", "Sets the URL to pipe standard error to.", &config.Run.StderrURL)
	runCmd.AddOptEnvString("stdout-output", 'o', "{OUTPUT|FILE}", "Sets the output destination for standard output.", &config.Run.StdoutOutput, flagx.WithDefaults(stringsx.SelectString(pex.IsGUI(), "null", "stdout")))
	runCmd.AddOptEnvString("stderr-output", 'e', "{OUTPUT|FILE}", "Sets the output destination for standard error.", &config.Run.StderrOutput, flagx.WithDefaults(stringsx.SelectString(pex.IsGUI(), "null", "stderr")))
	runCmd.AddOpt("debug", 'd', "", "Sets the http://localhost:8888/ URL to pipe standard streams to.", flagx.Join(flagx.Args(flagx.String(&config.Run.StdoutURL), "http://localhost:8888/"), flagx.Args(flagx.String(&config.Run.StderrURL), "http://localhost:8888/")))
	runCmd.AddOptBool("persistent", 'p', "", "Sets the command to run persistently.", &config.Run.Persistent, flagx.WithArgs("true"))
	runCmd.AddOptEnvString("broker-url", 'b', "URL", "Sets the URL to the broker.", &config.Run.BrokerURL, flagx.WithDefaults("http://localhost:8668/api/v1/processes"))
	runCmd.AddOptEnvDuration("broker-wait-timeout", 't', "DURATION", "Sets the wait timeout for the broker.", &config.Run.BrokerWaitTimeout, flagx.WithDefaults("10s"))
	runCmd.AddOptString("broker-client-name", 'n', "NAME", "Sets the client name for the broker.", &config.Run.BrokerClientName)
	runCmd.SetDefaultHandlerParams(stringsx.SelectString(pex.IsGUI(), "COMMAND [ARG ...]", "[COMMAND [ARG ...]]"), flagx.SelectValue(pex.IsGUI(), flagx.Join(flagx.String(&config.Run.CommandName), flagx.StringSlice(&config.Run.CommandArgs)), flagx.Optional(flagx.Join(flagx.String(&config.Run.CommandName), flagx.StringSlice(&config.Run.CommandArgs)))))
	runCmd.AddParam("URL", "The target URL. Example: http://localhost:8888/")
	runCmd.AddParam("COMMAND", "The name of the command.")
	runCmd.AddParam("ARG", "The arguments to the command.")
	runCmd.AddParam("OUTPUT", "The output destination. One of: stdout, stderr, null.")
	runCmd.AddParam("FILE", "The file path. Example: \"file.txt\".")
	runCmd.AddParam("BOOL", "The boolean value. One of: true, false.")
	runCmd.AddParam("DURATION", "The duration value. Example: 10s.")
	runCmd.AddParam("NAME", "The name value. Example: name.")

	debugCmd := flagx.AddCmd("debug")
	debugCmd.AddHelp(iox.SelectWriter(pex.IsGUI(), pex.GUIOutput(pex.IconInformation(), appname), os.Stderr), true)
	debugCmd.SetShortUsage("Runs the debug HTTP server.")
	debugCmd.AddOptString("address", 'a', "ADDRESS", "Sets the address to bind the debug HTTP server to.", &config.Debug.Address, flagx.WithDefaults("localhost:8888"))
	debugCmd.AddOptString("stdout-output", 'o', "{OUTPUT|FILE}", "Sets the output destination for received messages with stdout source.", &config.Debug.StdoutOutput, flagx.WithDefaults(stringsx.SelectString(pex.IsGUI(), "null", "stdout")))
	debugCmd.AddOptString("stderr-output", 'e', "{OUTPUT|FILE}", "Sets the output destination for received messages with stderr source.", &config.Debug.StderrOutput, flagx.WithDefaults(stringsx.SelectString(pex.IsGUI(), "null", "stderr")))
	debugCmd.AddParam("ADDRESS", "The local endpoint address. Example: localhost:8888")
	debugCmd.AddParam("OUTPUT", "The output destination. One of: stdout, stderr, null.")
	debugCmd.AddParam("FILE", "The file path. Example: \"file.txt\".")

	brokerCmd := flagx.AddCmd("broker")
	brokerCmd.AddPrefix("BROKER")
	brokerCmd.AddHelp(iox.SelectWriter(pex.IsGUI(), pex.GUIOutput(pex.IconInformation(), appname), os.Stderr), true)
	brokerCmd.SetShortUsage("Runs the broker HTTP server.")
	brokerCmd.AddOptEnvString("address", 'a', "ADDRESS", "Sets the address to bind the broker HTTP server to.", &config.Broker.Address, flagx.WithDefaults("localhost:8668"))
	brokerCmd.AddOptEnvString("stdout-output", 'o', "{OUTPUT|FILE}", "Sets the output destination for received messages.", &config.Broker.StdoutOutput, flagx.WithDefaults(stringsx.SelectString(pex.IsGUI(), "null", "stdout")))
	brokerCmd.AddOptEnvDuration("wait-timeout", 't', "DURATION", "Sets the wait timeout for the broker.", &config.Broker.WaitTimeout, flagx.WithDefaults("10s"))
	brokerCmd.AddParam("ADDRESS", "The local endpoint address. Example: localhost:8888")
	brokerCmd.AddParam("OUTPUT", "The output destination. One of: stdout, stderr, null.")
	brokerCmd.AddParam("DURATION", "The duration value. Example: 10s.")
	brokerCmd.AddParam("FILE", "The file path. Example: \"file.txt\".")

	listCmd := flagx.AddCmd("list")
	listCmd.AddHelp(iox.SelectWriter(pex.IsGUI(), pex.GUIOutput(pex.IconInformation(), appname), os.Stderr), true)
	listCmd.SetShortUsage("Lists the running processes.")
	listCmd.AddOptString("stdout-output", 'o', "{OUTPUT|FILE}", "Sets the output destination for standard output.", &config.List.StdoutOutput, flagx.WithDefaults(stringsx.SelectString(pex.IsGUI(), "null", "stdout")))
	listCmd.AddOptEnvString("broker-url", 'b', "URL", "Sets the URL to the broker.", &config.List.BrokerURL, flagx.WithDefaults("http://localhost:8668/api/v1/processes"))
	listCmd.AddParam("URL", "The target URL. Example: http://localhost:8888/")
	listCmd.AddParam("OUTPUT", "The output destination. One of: stdout, stderr, null.")
	listCmd.AddParam("FILE", "The file path. Example: \"file.txt\".")

	killCmd := flagx.AddCmd("kill")
	killCmd.AddHelp(iox.SelectWriter(pex.IsGUI(), pex.GUIOutput(pex.IconInformation(), appname), os.Stderr), true)
	killCmd.SetShortUsage("Kills the running process.")
	killCmd.SetDescription("Kills the running process by PID or pattern. Pattern never matches the parent process of broker.")
	killCmd.AddOptString("stdout-output", 'o', "{OUTPUT|FILE}", "Sets the output destination for standard output.", &config.Kill.StdoutOutput, flagx.WithDefaults(stringsx.SelectString(pex.IsGUI(), "null", "stdout")))
	killCmd.AddOptEnvString("broker-url", 'b', "URL", "Sets the URL to the broker.", &config.Kill.BrokerURL, flagx.WithDefaults("http://localhost:8668/api/v1/processes"))
	killCmd.SetDefaultHandlerParams("{PID|PATTERN}", flagx.String(&config.Kill.Pattern))
	killCmd.AddParam("PID", "The process ID. Example: 1234")
	killCmd.AddParam("PATTERN", "The process ID, client name or command name. Example: MYGROUP*.")
	killCmd.AddParam("URL", "The target URL. Example: http://localhost:8888/")
	killCmd.AddParam("OUTPUT", "The output destination. One of: stdout, stderr, null.")
	killCmd.AddParam("FILE", "The file path. Example: \"file.txt\".")

	flagx.Parse(os.Args[1:]...)
	return &config
}

package configs

import "time"

type StdhttpConfig struct {
	Run    StdhttpRunConfig
	Debug  StdhttpDebugConfig
	Broker StdhttpBrokerConfig
	List   StdhttpListConfig
	Kill   StdhttpKillConfig
}

type StdhttpRunConfig struct {
	StdoutURL    string
	StderrURL    string
	StdoutOutput string
	StderrOutput string

	CommandName string
	CommandArgs []string
	Persistent  bool

	BrokerURL         string
	BrokerClientName  string
	BrokerWaitTimeout time.Duration
}

type StdhttpDebugConfig struct {
	Address      string
	StdoutOutput string
	StderrOutput string
}

type StdhttpBrokerConfig struct {
	Address      string
	StdoutOutput string
	WaitTimeout  time.Duration
}

type StdhttpListConfig struct {
	BrokerURL    string
	StdoutOutput string
}

type StdhttpKillConfig struct {
	BrokerURL    string
	StdoutOutput string
	Pattern      string
}

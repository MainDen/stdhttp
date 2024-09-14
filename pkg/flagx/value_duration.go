package flagx

import (
	"errors"
	"time"
)

type valueDuration struct {
	value *time.Duration
}

func Duration(value *time.Duration) *valueDuration {
	return &valueDuration{value: value}
}

func (value *valueDuration) Parse(args ...string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("missing argument")
	}
	v, err := time.ParseDuration(args[0])
	if err != nil {
		return 0, err
	}
	*value.value = v
	return 1, nil
}

func (value *valueDuration) Format() string {
	return (*value.value).String()
}

func (fs *FlagSet) AddOptDuration(name string, alias rune, params string, usage string, pointer *time.Duration, wrappers ...Wrapper) *OptFlag {
	return fs.AddOpt(name, alias, params, usage, Duration(pointer), wrappers...)
}

func (fs *FlagSet) AddEnvDuration(name string, params string, usage string, pointer *time.Duration, wrappers ...Wrapper) *EnvFlag {
	return fs.AddEnv(name, params, usage, Duration(pointer), wrappers...)
}

func (fs *FlagSet) AddOptEnvDuration(name string, alias rune, params string, usage string, pointer *time.Duration, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return fs.AddOptEnv(name, alias, params, usage, Duration(pointer), wrappers...)
}

func AddOptDuration(name string, alias rune, params string, usage string, pointer *time.Duration, wrappers ...Wrapper) *OptFlag {
	return CommandLine.AddOpt(name, alias, params, usage, Duration(pointer), wrappers...)
}

func AddEnvDuration(name string, params string, usage string, pointer *time.Duration, wrappers ...Wrapper) *EnvFlag {
	return CommandLine.AddEnv(name, params, usage, Duration(pointer), wrappers...)
}

func AddOptEnvDuration(name string, alias rune, params string, usage string, pointer *time.Duration, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return CommandLine.AddOptEnv(name, alias, params, usage, Duration(pointer), wrappers...)
}

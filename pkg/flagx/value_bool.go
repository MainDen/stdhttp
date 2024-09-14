package flagx

import (
	"errors"
	"strconv"
)

type valueBool struct {
	value *bool
}

func Bool(value *bool) *valueBool {
	return &valueBool{value: value}
}

func (value *valueBool) Parse(args ...string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("missing argument")
	}
	v, err := strconv.ParseBool(args[0])
	if err != nil {
		return 0, err
	}
	*value.value = v
	return 1, nil
}

func (value *valueBool) Format() string {
	return strconv.FormatBool(*value.value)
}

func (fs *FlagSet) AddOptBool(name string, alias rune, params string, usage string, pointer *bool, wrappers ...Wrapper) *OptFlag {
	return fs.AddOpt(name, alias, params, usage, Bool(pointer), wrappers...)
}

func (fs *FlagSet) AddEnvBool(name string, params string, usage string, pointer *bool, wrappers ...Wrapper) *EnvFlag {
	return fs.AddEnv(name, params, usage, Bool(pointer), wrappers...)
}

func (fs *FlagSet) AddOptEnvBool(name string, alias rune, params string, usage string, pointer *bool, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return fs.AddOptEnv(name, alias, params, usage, Bool(pointer), wrappers...)
}

func AddOptBool(name string, alias rune, params string, usage string, pointer *bool, wrappers ...Wrapper) *OptFlag {
	return CommandLine.AddOpt(name, alias, params, usage, Bool(pointer), wrappers...)
}

func AddEnvBool(name string, params string, usage string, pointer *bool, wrappers ...Wrapper) *EnvFlag {
	return CommandLine.AddEnv(name, params, usage, Bool(pointer), wrappers...)
}

func AddOptEnvBool(name string, alias rune, params string, usage string, pointer *bool, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return CommandLine.AddOptEnv(name, alias, params, usage, Bool(pointer), wrappers...)
}

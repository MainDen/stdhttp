package flagx

import (
	"errors"
	"strconv"
)

type valueInt struct {
	pointer *int
}

func Int(pointer *int) *valueInt {
	return &valueInt{pointer: pointer}
}

func (value *valueInt) Parse(args ...string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("missing argument")
	}
	n, err := strconv.ParseInt(args[0], 0, 0)
	if err != nil {
		return 0, err
	}
	*value.pointer = int(n)
	return 1, nil
}

func (value *valueInt) Format() string {
	return strconv.Itoa(*value.pointer)
}

func (fs *FlagSet) AddOptInt(name string, alias rune, params string, usage string, pointer *int, wrappers ...Wrapper) *OptFlag {
	return fs.AddOpt(name, alias, params, usage, Int(pointer), wrappers...)
}

func (fs *FlagSet) AddEnvInt(name string, params string, usage string, pointer *int, wrappers ...Wrapper) *EnvFlag {
	return fs.AddEnv(name, params, usage, Int(pointer), wrappers...)
}

func (fs *FlagSet) AddOptEnvInt(name string, alias rune, params string, usage string, pointer *int, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return fs.AddOptEnv(name, alias, params, usage, Int(pointer), wrappers...)
}

func AddOptInt(name string, alias rune, params string, usage string, pointer *int, wrappers ...Wrapper) *OptFlag {
	return CommandLine.AddOpt(name, alias, params, usage, Int(pointer), wrappers...)
}

func AddEnvInt(name string, params string, usage string, pointer *int, wrappers ...Wrapper) *EnvFlag {
	return CommandLine.AddEnv(name, params, usage, Int(pointer), wrappers...)
}

func AddOptEnvInt(name string, alias rune, params string, usage string, pointer *int, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return CommandLine.AddOptEnv(name, alias, params, usage, Int(pointer), wrappers...)
}

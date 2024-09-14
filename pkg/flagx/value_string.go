package flagx

import "errors"

type valueString struct {
	pointer *string
}

func String(pointer *string) *valueString {
	return &valueString{pointer: pointer}
}

func (value *valueString) Parse(args ...string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("missing argument")
	}
	*value.pointer = args[0]
	return 1, nil
}

func (value *valueString) Format() string {
	return *value.pointer
}

func (fs *FlagSet) AddOptString(name string, alias rune, params string, usage string, pointer *string, wrappers ...Wrapper) *OptFlag {
	return fs.AddOpt(name, alias, params, usage, String(pointer), wrappers...)
}

func (fs *FlagSet) AddEnvString(name string, params string, usage string, pointer *string, wrappers ...Wrapper) *EnvFlag {
	return fs.AddEnv(name, params, usage, String(pointer), wrappers...)
}

func (fs *FlagSet) AddOptEnvString(name string, alias rune, params string, usage string, pointer *string, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return fs.AddOptEnv(name, alias, params, usage, String(pointer), wrappers...)
}

func AddOptString(name string, alias rune, params string, usage string, pointer *string, wrappers ...Wrapper) *OptFlag {
	return CommandLine.AddOpt(name, alias, params, usage, String(pointer), wrappers...)
}

func AddEnvString(name string, params string, usage string, pointer *string, wrappers ...Wrapper) *EnvFlag {
	return CommandLine.AddEnv(name, params, usage, String(pointer), wrappers...)
}

func AddOptEnvString(name string, alias rune, params string, usage string, pointer *string, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return CommandLine.AddOptEnv(name, alias, params, usage, String(pointer), wrappers...)
}

type valueStringSlice struct {
	pointer *[]string
}

func StringSlice(pointer *[]string) *valueStringSlice {
	return &valueStringSlice{pointer: pointer}
}

func (value *valueStringSlice) Parse(args ...string) (int, error) {
	*value.pointer = append([]string(nil), args...)
	return len(args), nil
}

func (value *valueStringSlice) Format() []string {
	return *value.pointer
}

func (fs *FlagSet) AddOptStringSlice(name string, alias rune, params string, usage string, pointer *[]string, wrappers ...Wrapper) *OptFlag {
	return fs.AddOpt(name, alias, params, usage, StringSlice(pointer), wrappers...)
}

func (fs *FlagSet) AddEnvStringSlice(name string, params string, usage string, pointer *[]string, wrappers ...Wrapper) *EnvFlag {
	return fs.AddEnv(name, params, usage, StringSlice(pointer), wrappers...)
}

func (fs *FlagSet) AddOptEnvStringSlice(name string, alias rune, params string, usage string, pointer *[]string, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return fs.AddOptEnv(name, alias, params, usage, StringSlice(pointer), wrappers...)
}

func AddOptStringSlice(name string, alias rune, params string, usage string, pointer *[]string, wrappers ...Wrapper) *OptFlag {
	return CommandLine.AddOpt(name, alias, params, usage, StringSlice(pointer), wrappers...)
}

func AddEnvStringSlice(name string, params string, usage string, pointer *[]string, wrappers ...Wrapper) *EnvFlag {
	return CommandLine.AddEnv(name, params, usage, StringSlice(pointer), wrappers...)
}

func AddOptEnvStringSlice(name string, alias rune, params string, usage string, pointer *[]string, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return CommandLine.AddOptEnv(name, alias, params, usage, StringSlice(pointer), wrappers...)
}

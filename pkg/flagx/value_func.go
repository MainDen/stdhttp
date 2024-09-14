package flagx

type valueFunc func(...string) (int, error)

func Func(f func(...string) (int, error)) valueFunc {
	return valueFunc(f)
}

func (f valueFunc) Parse(args ...string) (int, error) {
	return f(args...)
}

func (fs *FlagSet) AddOptFunc(name string, alias rune, params string, usage string, f func(...string) (int, error), wrappers ...Wrapper) *OptFlag {
	return fs.AddOpt(name, alias, params, usage, Func(f), wrappers...)
}

func (fs *FlagSet) AddEnvFunc(name string, params string, usage string, f func(...string) (int, error), wrappers ...Wrapper) *EnvFlag {
	return fs.AddEnv(name, params, usage, Func(f), wrappers...)
}

func (fs *FlagSet) AddOptEnvFunc(name string, alias rune, params string, usage string, f func(...string) (int, error), wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return fs.AddOptEnv(name, alias, params, usage, Func(f), wrappers...)
}

func AddOptFunc(name string, alias rune, params string, usage string, f func(...string) (int, error), wrappers ...Wrapper) *OptFlag {
	return CommandLine.AddOpt(name, alias, params, usage, Func(f), wrappers...)
}

func AddEnvFunc(name string, params string, usage string, f func(...string) (int, error), wrappers ...Wrapper) *EnvFlag {
	return CommandLine.AddEnv(name, params, usage, Func(f), wrappers...)
}

func AddOptEnvFunc(name string, alias rune, params string, usage string, f func(...string) (int, error), wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return CommandLine.AddOptEnv(name, alias, params, usage, Func(f), wrappers...)
}

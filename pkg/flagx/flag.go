package flagx

type EnvFlag struct {
	Name    string
	Params  string
	Usage   string
	Value   Value
	Default string
}

type OptFlag struct {
	Name     string
	Alias    rune
	Params   string
	Usage    string
	Value    Value
	Defaults []string
}

package flagx

import "io"

type valueHelp struct {
	fs     *FlagSet
	output io.Writer
	sealed bool
}

func Help(fs *FlagSet, output io.Writer, sealed bool) *valueHelp {
	return &valueHelp{fs: fs, output: output, sealed: sealed}
}

func (v *valueHelp) Parse(args ...string) (int, error) {
	if v.output != nil {
		output := v.fs.output
		v.fs.output = v.output
		defer func() { v.fs.output = output }()
	}
	if v.sealed {
		v.fs.Usage()
	} else {
		v.fs.Usage(args...)
	}
	return len(args), ErrHelp
}

func (fs *FlagSet) AddHelp(output io.Writer, sealed bool) *OptFlag {
	var params string
	if !sealed {
		params = "[Command ...]"
	}
	return fs.AddOpt("help", '?', params, "Prints command help.", Help(fs, output, sealed))
}

func AddHelp(output io.Writer, sealed bool) *OptFlag {
	return CommandLine.AddHelp(output, sealed)
}

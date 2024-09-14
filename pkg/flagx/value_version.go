package flagx

import (
	"fmt"
	"runtime"
	"strings"
)

type valueVersion struct {
	fs        *FlagSet
	version   string
	copyright string
	license   string
	url       string
}

func Version(fs *FlagSet, version string, copyright string, license string, url string) *valueVersion {
	return &valueVersion{fs: fs, version: version, copyright: copyright, license: license, url: url}
}

func (v *valueVersion) String() string {
	var builder strings.Builder
	_, _ = fmt.Fprintf(&builder, "%v version %v %v/%v (%v)\n", v.fs.name, v.version, runtime.GOOS, runtime.GOARCH, runtime.Version())
	if v.copyright != "" || v.license != "" || v.url != "" {
		_, _ = fmt.Fprintf(&builder, "\n")
	}
	if v.copyright != "" {
		_, _ = fmt.Fprintf(&builder, "Copyright (c) %v\n", v.copyright)
	}
	if v.license != "" {
		_, _ = fmt.Fprintf(&builder, "Licensed under the %v\n", v.license)
	}
	if v.url != "" {
		_, _ = fmt.Fprintf(&builder, "For more information, visit: %v\n", v.url)
	}
	return builder.String()
}

func (v *valueVersion) Parse(args ...string) (int, error) {
	if _, err := fmt.Fprint(v.fs.Output(), v.String()); err != nil {
		return 0, fmt.Errorf("failed to write version: %w", err)
	}
	return len(args), ErrHelp
}

func (fs *FlagSet) AddVersion(version string, copyright string, license string, url string) *OptFlag {
	if version == "" {
		version = "undefined"
	}
	return fs.AddOpt("version", 'V', "", "Prints command version.", Version(fs, version, copyright, license, url))
}

func AddVersion(version string, copyright string, license string, url string) *OptFlag {
	return CommandLine.AddVersion(version, copyright, license, url)
}

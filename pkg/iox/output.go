package iox

import (
	"errors"
	"io"
	"os"
	"strings"
)

var consoleModeEnabled = true

func IsConsoleModeEnabled() bool {
	return consoleModeEnabled
}

func SetConsoleModeEnabled(value bool) {
	consoleModeEnabled = value
}

func Output(output string) (io.Writer, error) {
	switch strings.ToLower(output) {
	case "stdout":
		if !consoleModeEnabled {
			return nil, errors.New("stdout is available only in console mode")
		}
		return NewOnlyWriter(os.Stdout), nil
	case "stderr":
		if !consoleModeEnabled {
			return nil, errors.New("stderr is available only in console mode")
		}
		return NewOnlyWriter(os.Stderr), nil
	case "null":
		return io.Discard, nil
	default:
		return os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	}
}

func SelectWriter(predicate bool, trueValue io.Writer, falseValue io.Writer) io.Writer {
	if predicate {
		return trueValue
	}
	return falseValue
}

type onlyWriter struct {
	w io.Writer
}

func (w *onlyWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func NewOnlyWriter(w io.Writer) io.Writer {
	return &onlyWriter{w}
}

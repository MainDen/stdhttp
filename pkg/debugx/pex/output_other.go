//go:build !windows

package pex

import (
	"io"
)

func IconError() int {
	return 0
}

func IconInformation() int {
	return 0
}

func GUIOutput(icon int, title string) io.Writer {
	return io.Discard
}

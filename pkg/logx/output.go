package logx

import (
	"errors"
	"io"
	"os"
	"sync/atomic"

	"github.com/mainden/stdhttp/pkg/iox"
)

var defaultOutput atomicWriter

func init() {
	SetDefaultOutput(iox.NewOnlyWriter(os.Stdout))
}

func ParseOutput(args ...string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("missing output")
	}
	output, err := iox.Output(args[0])
	if err != nil {
		return 0, err
	}
	SetDefaultOutput(output)
	return 1, nil
}

func DefaultOutput() io.Writer {
	return &defaultOutput
}

func SetDefaultOutput(output io.Writer) {
	defaultOutput.Store(output)
}

type atomicWriter struct {
	w atomic.Pointer[io.Writer]
}

func (output *atomicWriter) Load() io.Writer {
	if w := output.w.Load(); w != nil {
		return *w
	}
	return nil
}

func (output *atomicWriter) Store(w io.Writer) {
	output.w.Store(&w)
}

func (w *atomicWriter) Write(p []byte) (int, error) {
	return w.Load().Write(p)
}

func Close() error {
	return iox.Close(DefaultOutput())
}

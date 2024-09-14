package iox

import (
	"context"
	"io"

	"github.com/mainden/stdhttp/pkg/runx"
)

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func NewContextReader(ctx context.Context, r io.Reader) *contextReader {
	return &contextReader{
		ctx: ctx,
		r:   r,
	}
}

func (r *contextReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	case <-runx.Async(func() {
		n, err = r.r.Read(p)
	}):
		return n, err
	}
}

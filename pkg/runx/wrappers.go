package runx

import (
	"context"
	"time"
)

func Async(fns ...func()) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		for _, fn := range fns {
			fn()
		}
	}()
	return ch
}

func Await(ch <-chan struct{}, fns ...func()) {
	<-ch
	for _, fn := range fns {
		fn()
	}
}

func AwaitWithTimeout(timeout time.Duration, ch <-chan struct{}, fns ...func()) {
	select {
	case <-time.After(timeout):
	case <-ch:
	}
	for _, fn := range fns {
		fn()
	}
}

func AwaitDone(ctx context.Context, fns ...func()) {
	Await(ctx.Done(), fns...)
}

func AwaitDoneWithTimeout(ctx context.Context, timeout time.Duration, fns ...func()) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	AwaitDone(ctx, fns...)
}

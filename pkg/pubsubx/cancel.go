package pubsubx

import "context"

type cancelContextKey struct{}

func WithCancel(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	return context.WithValue(ctx, cancelContextKey{}, cancel)
}

func Cancel(ctx context.Context) {
	if cancel, ok := ctx.Value(cancelContextKey{}).(context.CancelFunc); ok {
		cancel()
	}
}

package logx

import (
	"context"
)

var defaultRegistry contextRegistry

type contextRegistry struct {
	args    []string
	getters map[string]func(context.Context) any
}

func GetContextArgs(ctx context.Context, args ...string) []any {
	return defaultRegistry.GetContextArgs(ctx, args...)
}

func (registery contextRegistry) GetContextArgs(ctx context.Context, args ...string) []any {
	if len(args) == 0 {
		args = registery.args
	}
	var result []any
	for _, arg := range args {
		if getter := registery.getters[arg]; getter != nil {
			value := getter(ctx)
			if value != nil {
				result = append(result, arg, value)
			}
		}
	}
	return result
}

type getterAnyFunc[T any] (func(context.Context) T)

func (getter getterAnyFunc[T]) Get(ctx context.Context) any {
	value := getter(ctx)
	return value
}

type getterNonZeroFunc[T comparable] (func(context.Context) T)

func (getter getterNonZeroFunc[T]) Get(ctx context.Context) any {
	value := getter(ctx)
	var zero T
	if value == zero {
		return nil
	}
	return value
}

func registerContextArg(arg string, getter func(context.Context) any) {
	if defaultRegistry.getters[arg] != nil {
		panic("duplicate context arg: " + arg)
	}
	if defaultRegistry.getters == nil {
		defaultRegistry.getters = make(map[string]func(context.Context) any)
	}
	defaultRegistry.args = append(defaultRegistry.args, arg)
	defaultRegistry.getters[arg] = getter
}

func UnregisterContextArg(arg string) {
	delete(defaultRegistry.getters, arg)
}

func RegisterContextAnyArg[T any](arg string, getter func(context.Context) T) {
	registerContextArg(arg, getterAnyFunc[T](getter).Get)
}

func RegisterContextNonZeroArg[T comparable](arg string, getter func(context.Context) T) {
	registerContextArg(arg, getterNonZeroFunc[T](getter).Get)
}

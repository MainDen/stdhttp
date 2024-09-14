package logx

import (
	"context"
	"encoding/base32"
	"encoding/binary"
	"maps"
	"math/rand"
	"strings"
	"time"
)

func init() {
	RegisterContextNonZeroArg("context_name", GetName)
	RegisterContextNonZeroArg("event_base_id", GetEventBaseId)
	RegisterContextNonZeroArg("event_last_id", GetEventLastId)
	RegisterContextNonZeroArg("event_id", GetEventId)
	RegisterContextAnyArg("context_values", GetValuesAny)
}

var eventRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type nameContextKey struct{}

type eventBaseContextKey struct{}

type eventLastContextKey struct{}

type eventContextKey struct{}

type valuesContextKey struct{}

func SetName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, nameContextKey{}, name)
}

func WithName(ctx context.Context, name string) context.Context {
	if base, ok := ctx.Value(nameContextKey{}).(string); ok {
		return context.WithValue(ctx, nameContextKey{}, base+"."+name)
	}
	return context.WithValue(ctx, nameContextKey{}, name)
}

func GetName(ctx context.Context) string {
	if name, ok := ctx.Value(nameContextKey{}).(string); ok {
		return name
	}
	return ""
}

func MakeEventId(event string) string {
	return event + ":" + strings.ToLower(base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(binary.LittleEndian.AppendUint64(make([]byte, 0, 8), eventRand.Uint64())))
}

func SetEventId(ctx context.Context, eventId string) context.Context {
	ctx = context.WithValue(ctx, eventBaseContextKey{}, eventId)
	ctx = context.WithValue(ctx, eventLastContextKey{}, "")
	return context.WithValue(ctx, eventContextKey{}, eventId)
}

func WithEventId(ctx context.Context, eventId string) context.Context {
	if _, ok := ctx.Value(eventBaseContextKey{}).(string); !ok {
		ctx = context.WithValue(ctx, eventBaseContextKey{}, eventId)
	}
	if eventLastId, ok := ctx.Value(eventContextKey{}).(string); ok {
		ctx = context.WithValue(ctx, eventLastContextKey{}, eventLastId)
	}
	return context.WithValue(ctx, eventContextKey{}, eventId)
}

func SetEvent(ctx context.Context, event string) context.Context {
	return SetEventId(ctx, MakeEventId(event))
}

func WithEvent(ctx context.Context, event string) context.Context {
	return WithEventId(ctx, MakeEventId(event))
}

func GetEventId(ctx context.Context) string {
	if name, ok := ctx.Value(eventContextKey{}).(string); ok {
		return name
	}
	return ""
}

func GetEventLastId(ctx context.Context) string {
	if name, ok := ctx.Value(eventLastContextKey{}).(string); ok {
		return name
	}
	return ""
}

func GetEventBaseId(ctx context.Context) string {
	if name, ok := ctx.Value(eventBaseContextKey{}).(string); ok {
		return name
	}
	return ""
}

func WithValue(ctx context.Context, key string, value any) context.Context {
	if m, ok := ctx.Value(valuesContextKey{}).(map[string]any); ok {
		m = maps.Clone(m)
		m[key] = value
		return context.WithValue(ctx, valuesContextKey{}, m)
	}
	return context.WithValue(ctx, valuesContextKey{}, map[string]any{key: value})
}

func GetValue(ctx context.Context, key string) any {
	if m, ok := ctx.Value(valuesContextKey{}).(map[string]any); ok {
		return m[key]
	}
	return nil
}

func GetValues(ctx context.Context) map[string]any {
	if m, ok := ctx.Value(valuesContextKey{}).(map[string]any); ok {
		return m
	}
	return nil
}

func GetValuesAny(ctx context.Context) any {
	if values := GetValues(ctx); values != nil {
		return values
	}
	return nil
}

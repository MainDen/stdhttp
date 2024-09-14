package logx

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"strings"
)

const (
	LevelDebug  slog.Level = slog.LevelDebug
	LevelInfo   slog.Level = slog.LevelInfo
	LevelWarn   slog.Level = slog.LevelWarn
	LevelError  slog.Level = slog.LevelError
	LevelSilent slog.Level = math.MaxInt
)

var defaultLevel slog.LevelVar

func ParseLevel(args ...string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("missing level")
	}
	switch strings.ToLower(args[0]) {
	case "debug":
		SetDefaultLevel(LevelDebug)
	case "info":
		SetDefaultLevel(LevelInfo)
	case "warn":
		SetDefaultLevel(LevelWarn)
	case "error":
		SetDefaultLevel(LevelError)
	case "silent":
		SetDefaultLevel(LevelSilent)
	default:
		return 0, errors.New("invalid level")
	}
	return 1, nil
}

func SetDefaultLevel(level slog.Level) {
	defaultLevel.Set(level)
}

func DefaultLevel() slog.Level {
	return defaultLevel.Level()
}

func DefaultLeveler() slog.Leveler {
	return &defaultLevel
}

type levelContextKey struct{}

func WithLevel(ctx context.Context, level slog.Level) context.Context {
	return context.WithValue(ctx, levelContextKey{}, level)
}

func GetLevel(ctx context.Context) slog.Level {
	if level, ok := ctx.Value(levelContextKey{}).(slog.Level); ok {
		return level
	}
	return DefaultLevel()
}

type levelHandler struct {
	handler slog.Handler
}

func WrapLevelHandler(handler slog.Handler) slog.Handler {
	return &levelHandler{handler: handler}
}

func (h *levelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= GetLevel(ctx)
}

func (h *levelHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.handler.Handle(ctx, record)
}

func (h *levelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return WrapLevelHandler(h.handler.WithAttrs(attrs))
}

func (h *levelHandler) WithGroup(name string) slog.Handler {
	return WrapLevelHandler(h.handler.WithGroup(name))
}

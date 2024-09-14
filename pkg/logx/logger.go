package logx

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
)

var defaultLogger atomic.Pointer[Logger]

func init() {
	SetDefault(NewJsonLogger())
}

func ParseFormat(args ...string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("missing format")
	}
	switch strings.ToLower(args[0]) {
	case "json":
		SetDefault(NewJsonLogger())
	case "text":
		SetDefault(NewTextLogger())
	default:
		return 0, errors.New("invalid format")
	}
	return 1, nil
}

func SetDefault(logger Logger) {
	defaultLogger.Store(&logger)
}

func Default() Logger {
	return *defaultLogger.Load()
}

type loggerContextKey struct{}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

func GetLogger(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey{}).(Logger); ok {
		return logger
	}
	return Default()
}

type Logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type contextLogger struct {
	logger Logger
	args   []string
}

func WrapLogger(logger Logger, args ...string) Logger {
	return &contextLogger{logger: logger, args: args}
}

func NewTextLogger(args ...string) Logger {
	return WrapLogger(slog.New(WrapLevelHandler(slog.NewTextHandler(DefaultOutput(), nil))), args...)
}

func NewJsonLogger(args ...string) Logger {
	return WrapLogger(slog.New(WrapLevelHandler(slog.NewJSONHandler(DefaultOutput(), nil))), args...)
}

func (l *contextLogger) ContextArgs(ctx context.Context) []any {
	return GetContextArgs(ctx, l.args...)
}

func (l *contextLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, append(l.ContextArgs(ctx), args...)...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	GetLogger(ctx).DebugContext(ctx, msg, args...)
}

func (l *contextLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, append(l.ContextArgs(ctx), args...)...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	GetLogger(ctx).InfoContext(ctx, msg, args...)
}

func (l *contextLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, append(l.ContextArgs(ctx), args...)...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	GetLogger(ctx).WarnContext(ctx, msg, args...)
}

func (l *contextLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, append(l.ContextArgs(ctx), args...)...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	GetLogger(ctx).ErrorContext(ctx, msg, args...)
}

func (l *contextLogger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

func FatalContext(ctx context.Context, msg string, args ...any) {
	if logger, ok := GetLogger(ctx).(interface {
		FatalContext(context.Context, string, ...any)
	}); ok {
		logger.FatalContext(ctx, msg, args...)
	}
	GetLogger(ctx).ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

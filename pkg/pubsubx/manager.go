package pubsubx

import (
	"context"
	"errors"
	"fmt"

	"github.com/mainden/stdhttp/pkg/logx"
	"github.com/mainden/stdhttp/pkg/runtimex"
)

var Default Manager = NewManager()

var ErrUnexpectedMessageType = errors.New("unexpected message type")

type Handler interface {
	Handle(ctx context.Context, message interface{}) error
}

type HandlerFunc func(ctx context.Context, message interface{}) error

func (f HandlerFunc) Handle(ctx context.Context, message interface{}) error {
	return f(ctx, message)
}

type Manager interface {
	Subscribe(ctx context.Context, topic string, handler Handler)
	Publish(ctx context.Context, topic string, message interface{}) error
}

type managerContextKey struct{}

func WithManager(ctx context.Context, manager Manager) context.Context {
	return context.WithValue(ctx, managerContextKey{}, manager)
}

func GetManager(ctx context.Context) Manager {
	if manager, ok := ctx.Value(managerContextKey{}).(Manager); ok {
		return manager
	}
	return Default
}

func Subscribe(ctx context.Context, topic string, handler Handler) {
	GetManager(ctx).Subscribe(ctx, topic, handler)
}

func Publish(ctx context.Context, topic string, message interface{}) error {
	return GetManager(ctx).Publish(ctx, topic, message)
}

type manager struct {
	handlers map[string][]Handler
}

func NewManager() *manager {
	return &manager{
		handlers: make(map[string][]Handler),
	}
}

func GetHandlerName(handler Handler) string {
	if _, ok := handler.(HandlerFunc); ok {
		return runtimex.GetFuncName(handler)
	}
	return fmt.Sprintf("%T", handler)
}

func (m *manager) Subscribe(ctx context.Context, topic string, handler Handler) {
	ctx = logx.WithName(ctx, "pubsub")
	m.handlers[topic] = append(m.handlers[topic], handler)
	logx.DebugContext(ctx, "Subscribed to topic", "topic", topic, "handler", GetHandlerName(handler))
}

func (m *manager) Publish(ctx context.Context, topic string, message interface{}) error {
	ctx = logx.WithName(ctx, "pubsub")
	ctx = logx.WithEvent(ctx, "pubsub:"+topic)
	logx.DebugContext(ctx, "Message", "topic", topic, "message", message)
	handlers, ok := m.handlers[topic]
	if !ok {
		logx.DebugContext(ctx, "No handlers", "topic", topic)
		return nil
	}
	var errs error
	for _, handler := range handlers {
		handlerName := GetHandlerName(handler)
		logx.DebugContext(ctx, "Publishing", "topic", topic, "handler", handlerName)
		if err := handler.Handle(ctx, message); err != nil {
			logx.DebugContext(ctx, "Failed to publish", "topic", topic, "handler", handlerName, "error", err)
			errs = errors.Join(errs, fmt.Errorf("failed to publish to topic '%s' using handler '%s': %w", topic, handlerName, err))
			continue
		}
		logx.DebugContext(ctx, "Published", "topic", topic, "handler", handlerName)
	}
	return errs
}

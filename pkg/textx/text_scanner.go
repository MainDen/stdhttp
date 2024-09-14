package textx

import (
	"bufio"
	"context"
	"io"

	"github.com/mainden/stdhttp/pkg/logx"
	"github.com/mainden/stdhttp/pkg/pubsubx"
)

type textScanner struct {
	reader io.Reader
	topic  string
}

func NewTextScanner(reader io.Reader, topic string) *textScanner {
	return &textScanner{
		reader: reader,
		topic:  topic,
	}
}

func (s *textScanner) Run(ctx context.Context) {
	ctx = logx.WithName(ctx, "text_scanner")
	scanner := bufio.NewScanner(s.reader)
	logx.DebugContext(ctx, "Running", "topic", s.topic)

	for scanner.Scan() {
		ctx := logx.SetEvent(ctx, "text_scanner")
		if err := pubsubx.Publish(ctx, s.topic, scanner.Text()); err != nil {
			logx.DebugContext(ctx, "Failed to publish message", "topic", s.topic, "error", err)
		}
	}
	if err := scanner.Err(); err != nil {
		logx.DebugContext(ctx, "Failed to scan text", "topic", s.topic, "error", err)
	}
	logx.DebugContext(ctx, "Stopped", "topic", s.topic)
}

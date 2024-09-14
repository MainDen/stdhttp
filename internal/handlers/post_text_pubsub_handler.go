package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mainden/stdhttp/internal/models"
	"github.com/mainden/stdhttp/pkg/httpx"
	"github.com/mainden/stdhttp/pkg/logx"
	"github.com/mainden/stdhttp/pkg/pubsubx"
)

type postTextPubsubHandler struct {
	url    string
	source string
}

func NewPostTextPubsubHandler(url string, source string) *postTextPubsubHandler {
	return &postTextPubsubHandler{
		url:    url,
		source: source,
	}
}

func (h *postTextPubsubHandler) Handle(ctx context.Context, message interface{}) (err error) {
	ctx = logx.WithName(ctx, "post_text_pubsub_handler")
	var body models.PostTextBody
	switch message := message.(type) {
	case string:
		body = models.MakePostTextBody(h.source, message)
	case []string:
		body = models.MakePostTextBody(h.source, message...)
	default:
		return fmt.Errorf("%w (%T)", pubsubx.ErrUnexpectedMessageType, message)
	}
	if len(body.Items) == 0 {
		logx.WarnContext(ctx, "No items to post")
		return nil
	}
	var resp *http.Response
	if resp, err = httpx.DoJson(ctx, http.MethodPost, h.url, &body); err != nil {
		return err
	}
	if err = httpx.AsNothing(resp.Body); err != nil {
		return err
	}
	return nil
}

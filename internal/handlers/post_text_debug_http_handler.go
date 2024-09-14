package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mainden/stdhttp/internal/models"
	"github.com/mainden/stdhttp/pkg/httpx"
	"github.com/mainden/stdhttp/pkg/logx"
)

type postTextDebugHttpHandler struct {
	stdout io.Writer
	stderr io.Writer
}

func NewPostTextDebugHttpHandler(stdout io.Writer, stderr io.Writer) *postTextDebugHttpHandler {
	return &postTextDebugHttpHandler{
		stdout: stdout,
		stderr: stderr,
	}
}

func (h *postTextDebugHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := logx.WithName(r.Context(), "post_text_debug_http_handler")
	var data []byte
	var err error
	if err = httpx.AsData(r.Body, &data); err != nil {
		logx.ErrorContext(ctx, "Failed to read body", "remote_address", r.RemoteAddr, "method", r.Method, "url", r.URL.String(), "headers", r.Header, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(data) == 0 {
		logx.ErrorContext(ctx, "Request received without body", "remote_address", r.RemoteAddr, "method", r.Method, "url", r.URL.String(), "headers", r.Header)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var body models.PostTextBody
	if err = json.Unmarshal(data, &body); err != nil {
		logx.ErrorContext(ctx, "Request received with invalid body format", "remote_address", r.RemoteAddr, "method", r.Method, "url", r.URL.String(), "headers", r.Header, "body", string(data))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logx.DebugContext(ctx, "Request received", "remote_address", r.RemoteAddr, "method", r.Method, "url", r.URL.String(), "headers", r.Header, "body", body)
	for _, item := range body.Items {
		switch item.Source {
		case "stdout":
			h.stdout.Write(append([]byte(item.Message), '\n'))
		case "stderr":
			h.stderr.Write(append([]byte(item.Message), '\n'))
		default:
			logx.WarnContext(ctx, "Unknown source", "source", item.Source)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

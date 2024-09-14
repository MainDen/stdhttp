package httpx

import (
	"net/http"

	"github.com/mainden/stdhttp/pkg/logx"
)

func HandleEvent(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if eventBaseId := r.Header.Get("X-Event-Base-Id"); eventBaseId != "" {
			ctx = logx.SetEventId(ctx, eventBaseId)
		}
		if eventId := r.Header.Get("X-Event-Id"); eventId != "" {
			ctx = logx.WithEventId(ctx, eventId)
		}
		ctx = logx.WithEvent(ctx, "http_handler")
		w.Header().Set("X-Event-Base-Id", logx.GetEventBaseId(ctx))
		w.Header().Set("X-Event-Id", logx.GetEventId(ctx))
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

func HandleLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logx.DebugContext(r.Context(), "HTTP request received", "method", r.Method, "url", r.URL)
		bw := NewStatisticResponseWriter(w)
		defer func() {
			if bw.Err() != nil {
				logx.DebugContext(r.Context(), "HTTP response failed", "method", r.Method, "url", r.URL, "status_code", bw.StatusCode(), "elapsed_time", bw.ElapsedTime(), "size", bw.Size(), "error", bw.Err())
				return
			}
			logx.DebugContext(r.Context(), "HTTP response sent", "method", r.Method, "url", r.URL, "status_code", bw.StatusCode(), "elapsed_time", bw.ElapsedTime(), "size", bw.Size())
		}()
		handler.ServeHTTP(bw, r)
	})
}

package httpx

import (
	"errors"
	"net/http"
	"time"
)

type statisticResponseWriter struct {
	wrapped    http.ResponseWriter
	statusCode int
	size       int
	timestamp  time.Time
	err        error
}

func NewStatisticResponseWriter(wrapped http.ResponseWriter) *statisticResponseWriter {
	return &statisticResponseWriter{
		wrapped:   wrapped,
		timestamp: time.Now(),
	}
}

func (w *statisticResponseWriter) Header() http.Header {
	return w.wrapped.Header()
}

func (w *statisticResponseWriter) Write(data []byte) (int, error) {
	n, err := w.wrapped.Write(data)
	w.size += n
	w.err = errors.Join(w.err, err)
	return n, err
}

func (w *statisticResponseWriter) WriteHeader(statusCode int) {
	w.wrapped.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *statisticResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *statisticResponseWriter) Size() int {
	return w.size
}

func (w *statisticResponseWriter) ElapsedTime() time.Duration {
	return time.Since(w.timestamp)
}

func (w *statisticResponseWriter) Err() error {
	return w.err
}

package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mainden/stdhttp/pkg/errorsx"
	"github.com/mainden/stdhttp/pkg/logx"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

func MakeErrorUnexpectedStatusCode(statusCode int) error {
	return fmt.Errorf("%w (%d)", ErrUnexpectedStatusCode, statusCode)
}

var Default HttpClient = WrapHttpClient(http.DefaultClient, WithLogger(), WithEvent())

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClientFunc func(req *http.Request) (*http.Response, error)

func (client HttpClientFunc) Do(req *http.Request) (*http.Response, error) {
	return client(req)
}

type httpClientContextKey struct{}

func WithHttpClient(ctx context.Context, client HttpClient) context.Context {
	return context.WithValue(ctx, httpClientContextKey{}, client)
}

func GetHttpClient(ctx context.Context) HttpClient {
	if client, ok := ctx.Value(httpClientContextKey{}).(HttpClient); ok {
		return client
	}
	return Default
}

func Do(req *http.Request) (*http.Response, error) {
	return GetHttpClient(req.Context()).Do(req)
}

func DoReader(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	return resp, nil
}

func DoJson(ctx context.Context, method string, url string, body any) (*http.Response, error) {
	bodyData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	return resp, nil
}

func dispose(body io.ReadCloser) func() error {
	return func() error {
		if body == nil {
			return nil
		}
		if err := body.Close(); err != nil {
			return fmt.Errorf("failed to close body: %w", err)
		}
		return nil
	}
}

func AsNothing(body io.ReadCloser) (err error) {
	defer errorsx.Dispose(&err, dispose(body))
	return nil
}

func AsData(body io.ReadCloser, dst *[]byte) (err error) {
	defer errorsx.Dispose(&err, dispose(body))
	*dst, err = io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	return nil
}

func WriteData(w http.ResponseWriter, status int, src []byte) (err error) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(status)
	if _, err := w.Write(src); err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}
	return nil
}

func AsJson(body io.ReadCloser, dst any) (err error) {
	defer errorsx.Dispose(&err, dispose(body))
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("failed to decode body: %w", err)
	}
	return nil
}

func WriteJson(w http.ResponseWriter, status int, src any) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(src); err != nil {
		return fmt.Errorf("failed to encode body: %w", err)
	}
	return nil
}

type Wrapper interface {
	Wrap(client HttpClient) HttpClient
}

type WrapperFunc func(client HttpClient) HttpClient

func (wrapper WrapperFunc) Wrap(client HttpClient) HttpClient {
	return wrapper(client)
}

func WrapHttpClient(client HttpClient, wrappers ...Wrapper) HttpClient {
	for _, wrapper := range wrappers {
		client = wrapper.Wrap(client)
	}
	return client
}

func WithUserAgent(userAgent string) Wrapper {
	return WrapperFunc(func(client HttpClient) HttpClient {
		return HttpClientFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Set("User-Agent", userAgent)
			return client.Do(req)
		})
	})
}

func WithEvent() Wrapper {
	return WrapperFunc(func(client HttpClient) HttpClient {
		return HttpClientFunc(func(req *http.Request) (*http.Response, error) {
			req = req.WithContext(logx.WithEvent(req.Context(), "http_request"))
			req.Header = req.Header.Clone()
			req.Header.Set("X-Event-Base-Id", logx.GetEventBaseId(req.Context()))
			req.Header.Set("X-Event-Id", logx.GetEventId(req.Context()))
			return client.Do(req)
		})
	})
}

func WithLogger() Wrapper {
	return WrapperFunc(func(client HttpClient) HttpClient {
		return HttpClientFunc(func(req *http.Request) (*http.Response, error) {
			req = req.WithContext(logx.WithName(req.Context(), "http_client"))
			logx.DebugContext(req.Context(), "HTTP request", "method", req.Method, "url", req.URL.String())
			resp, err := client.Do(req)
			if err != nil {
				logx.DebugContext(req.Context(), "HTTP request failed", "method", req.Method, "url", req.URL.String(), "error", err.Error())
			} else {
				logx.DebugContext(req.Context(), "HTTP response", "method", req.Method, "url", req.URL.String(), "status_code", resp.StatusCode)
			}
			return resp, err
		})
	})
}

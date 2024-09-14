package clients

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/mainden/stdhttp/internal/models"
	"github.com/mainden/stdhttp/pkg/httpx"
	"github.com/mainden/stdhttp/pkg/logx"
	"github.com/mainden/stdhttp/pkg/pubsubx"
	"github.com/mainden/stdhttp/pkg/runx"
)

type processesBrokerHttpClient struct {
	waitTimeout time.Duration
	url         string
}

func NewProcessesBrokerHttpClient(url string, waitTimeout time.Duration) *processesBrokerHttpClient {
	if waitTimeout <= 0 {
		waitTimeout = 10 * time.Second
	}
	return &processesBrokerHttpClient{
		url:         url,
		waitTimeout: waitTimeout,
	}
}

func (client *processesBrokerHttpClient) WaitTimeout() time.Duration {
	return client.waitTimeout
}

func (client *processesBrokerHttpClient) Register(ctx context.Context, process models.ProcessModel) (err error) {
	var resp *http.Response
	if resp, err = httpx.DoJson(ctx, http.MethodPost, client.url, models.MakeProcessesBodyItem(process)); err != nil {
		return err
	}
	if err = httpx.AsNothing(resp.Body); err != nil {
		return err
	}
	if resp.StatusCode == http.StatusCreated {
		return nil
	}
	if resp.StatusCode == http.StatusConflict {
		return models.ErrProcessExists
	}
	return httpx.MakeErrorUnexpectedStatusCode(resp.StatusCode)
}

func (client *processesBrokerHttpClient) Kill(ctx context.Context, pid int) (err error) {
	var resp *http.Response
	if resp, err = httpx.DoJson(ctx, http.MethodDelete, client.url+"/"+models.PidPathItem(pid), nil); err != nil {
		return err
	}
	if err = httpx.AsNothing(resp.Body); err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return models.ErrProcessNotFound
	}
	return httpx.MakeErrorUnexpectedStatusCode(resp.StatusCode)
}

func (client *processesBrokerHttpClient) KillMany(ctx context.Context, pattern string) (err error) {
	var resp *http.Response
	if resp, err = httpx.DoJson(ctx, http.MethodDelete, client.url+"?"+url.Values{"pattern": {pattern}}.Encode(), nil); err != nil {
		return err
	}
	if err = httpx.AsNothing(resp.Body); err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	return httpx.MakeErrorUnexpectedStatusCode(resp.StatusCode)
}

func (client *processesBrokerHttpClient) SendCommand(ctx context.Context, pid int, command string) (err error) {
	var resp *http.Response
	if resp, err = httpx.DoJson(ctx, http.MethodPut, client.url+"/"+models.PidPathItem(pid)+"/command", command); err != nil {
		return err
	}
	if err = httpx.AsNothing(resp.Body); err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return models.ErrProcessNotFound
	}
	if resp.StatusCode == http.StatusConflict {
		return models.ErrProcessBusy
	}
	return httpx.MakeErrorUnexpectedStatusCode(resp.StatusCode)
}

func (client *processesBrokerHttpClient) WaitCommand(ctx context.Context, pid int) (command string, err error) {
	ctx, cancel := context.WithTimeout(ctx, client.waitTimeout+time.Second)
	defer cancel()
	var resp *http.Response
	if resp, err = httpx.DoReader(ctx, http.MethodGet, client.url+"/"+models.PidPathItem(pid)+"/command", nil); err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusOK {
		if err = httpx.AsJson(resp.Body, &command); err != nil {
			return "", err
		}
		return command, nil
	}
	if err = httpx.AsNothing(resp.Body); err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusGone {
		return "", models.ErrProcessKilled
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", models.ErrProcessNotFound
	}
	if resp.StatusCode == http.StatusRequestTimeout {
		return "", models.ErrProcessWaitTimeout
	}
	return "", httpx.MakeErrorUnexpectedStatusCode(resp.StatusCode)
}

func (client *processesBrokerHttpClient) List(ctx context.Context) (processes models.ProcessModels, err error) {
	var resp *http.Response
	if resp, err = httpx.DoReader(ctx, http.MethodGet, client.url, nil); err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var body models.ProcessesBody
		if err = httpx.AsJson(resp.Body, &body); err != nil {
			return nil, err
		}
		return body.ProcessModels(), nil
	}
	if err = httpx.AsNothing(resp.Body); err != nil {
		return nil, err
	}
	return nil, httpx.MakeErrorUnexpectedStatusCode(resp.StatusCode)
}

func (client *processesBrokerHttpClient) Unregister(ctx context.Context, pid int) error {
	if err := client.Kill(ctx, pid); err != nil && !errors.Is(err, models.ErrProcessNotFound) {
		logx.DebugContext(ctx, "Failed to unregister process", "pid", pid, "error", err)
		return err
	}
	logx.DebugContext(ctx, "Process unregistered", "pid", pid)
	return nil
}

func (client *processesBrokerHttpClient) CommandLoop(ctx context.Context, process models.ProcessModel, topic string) {
	if err := client.Register(ctx, process); err != nil {
		logx.DebugContext(ctx, "Failed to register process", "pid", process.Pid, "error", err)
	} else {
		logx.DebugContext(ctx, "Registered process", "pid", process.Pid)
	}
	defer client.Unregister(context.Background(), process.Pid)
	for {
		command, err := client.WaitCommand(ctx, process.Pid)
		if err != nil && ctx.Err() != nil {
			logx.DebugContext(ctx, "Context canceled", "error", err)
			return
		}
		if err == nil && command == "" {
			logx.DebugContext(ctx, "Received empty command")
			continue
		}
		if err == nil {
			if err := pubsubx.Publish(ctx, topic, command); err != nil {
				logx.DebugContext(ctx, "Failed to publish command", "pid", process.Pid, "error", err)
				continue
			}
			logx.DebugContext(ctx, "Published command", "pid", process.Pid, "command", command)
			continue
		}
		if errors.Is(err, models.ErrProcessKilled) {
			logx.DebugContext(ctx, "Process killed", "pid", process.Pid)
			pubsubx.Cancel(ctx)
			return
		}
		if errors.Is(err, models.ErrProcessNotFound) {
			logx.DebugContext(ctx, "Process not found", "pid", process.Pid)
			if err := client.Register(ctx, process); err != nil {
				logx.DebugContext(ctx, "Failed to register process", "pid", process.Pid, "error", err)
				runx.AwaitDoneWithTimeout(ctx, client.waitTimeout)
				continue
			}
			logx.DebugContext(ctx, "Registered process", "pid", process.Pid)
			continue
		}
		if errors.Is(err, models.ErrProcessWaitTimeout) {
			logx.DebugContext(ctx, "Wait timeout", "pid", process.Pid)
			continue
		}
		logx.DebugContext(ctx, "Failed to wait command", "pid", process.Pid, "error", err)
		runx.AwaitDoneWithTimeout(ctx, client.waitTimeout)
	}
}

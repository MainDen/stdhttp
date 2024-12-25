package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/mainden/stdhttp/internal/models"
	"github.com/mainden/stdhttp/pkg/httpx"
	"github.com/mainden/stdhttp/pkg/slicesx"
)

type processesBroker interface {
	Register(ctx context.Context, process models.ProcessModel) (err error)
	Kill(ctx context.Context, pid int) (err error)
	SendCommand(ctx context.Context, pid int, command string) (err error)
	WaitCommand(ctx context.Context, pid int) (command string, err error)
	List(ctx context.Context) (processes models.ProcessModels, err error)
}

type processesBrokerHttpHandler struct {
	processesBroker processesBroker
	output          io.Writer
}

func NewProcessesBrokerHttpHandler(processesBroker processesBroker, output io.Writer) *processesBrokerHttpHandler {
	return &processesBrokerHttpHandler{
		processesBroker: processesBroker,
		output:          output,
	}
}

func (handler *processesBrokerHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch op := r.URL.Query().Get("op"); op {
	case "Register":
		handler.register(w, r)
	case "KillMany":
		handler.killMany(w, r)
	case "Kill":
		handler.kill(w, r)
	case "SendCommand":
		handler.sendCommand(w, r)
	case "WaitCommand":
		handler.waitCommand(w, r)
	case "List":
		handler.list(w, r)
	default:
		http.Error(w, fmt.Sprintf("unknown op '%v'", op), http.StatusBadRequest)
	}
}

func (handler *processesBrokerHttpHandler) register(w http.ResponseWriter, r *http.Request) {
	var body models.ProcessesBodyItem
	if err := httpx.AsJson(r.Body, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	process := body.ProcessModel()
	if err := handler.processesBroker.Register(r.Context(), process); err != nil {
		if errors.Is(err, models.ErrProcessExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(handler.output, "process '%v': registered\n", process.Pid)
	w.WriteHeader(http.StatusCreated)
}

func (handler *processesBrokerHttpHandler) killMany(w http.ResponseWriter, r *http.Request) {
	if err := httpx.AsNothing(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	processes, err := handler.processesBroker.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pattern := r.URL.Query().Get("pattern")
	processes = slicesx.Select(processes, func(process models.ProcessModel) bool {
		if pattern == "" {
			return false
		}
		matchedPid, _ := path.Match(pattern, strconv.Itoa(process.Pid))
		matchedClientName, _ := path.Match(pattern, process.ClientName)
		matchedCommandName, _ := path.Match(pattern, process.CommandName)
		return matchedPid || matchedClientName || matchedCommandName
	})

	for _, process := range processes {
		if os.Getppid() == process.Pid {
			continue
		}
		if err := handler.processesBroker.Kill(r.Context(), process.Pid); err != nil {
			if errors.Is(err, models.ErrProcessNotFound) {
				continue
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Fprintf(handler.output, "kill '%v': error\n", process.Pid)
			return
		}
		fmt.Fprintf(handler.output, "kill '%v': success\n", process.Pid)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (handler *processesBrokerHttpHandler) kill(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.Atoi(r.URL.Query().Get("pid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := httpx.AsNothing(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.processesBroker.Kill(r.Context(), pid); err != nil {
		if errors.Is(err, models.ErrProcessNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			fmt.Fprintf(handler.output, "kill '%v': process not found\n", pid)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(handler.output, "kill '%v': error\n", pid)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(handler.output, "kill '%v': success\n", pid)
}

func (handler *processesBrokerHttpHandler) sendCommand(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.Atoi(r.URL.Query().Get("pid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var command string
	if err := httpx.AsJson(r.Body, &command); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.processesBroker.SendCommand(r.Context(), pid, command); err != nil {
		if errors.Is(err, models.ErrProcessNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			fmt.Fprintf(handler.output, "send command '%v' to process '%v': process not found\n", command, pid)
			return
		}
		if errors.Is(err, models.ErrProcessBusy) {
			http.Error(w, err.Error(), http.StatusConflict)
			fmt.Fprintf(handler.output, "send command '%v' to process '%v': process busy\n", command, pid)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(handler.output, "send command '%v' to process '%v': unexpected error\n", command, pid)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(handler.output, "send command '%v' to process '%v': success\n", command, pid)
}

func (handler *processesBrokerHttpHandler) waitCommand(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.Atoi(r.URL.Query().Get("pid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := httpx.AsNothing(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	command, err := handler.processesBroker.WaitCommand(r.Context(), pid)
	if err != nil {
		if errors.Is(err, models.ErrProcessNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			fmt.Fprintf(handler.output, "process '%v': wait error: process not found\n", pid)
			return
		}
		if errors.Is(err, models.ErrProcessKilled) {
			http.Error(w, err.Error(), http.StatusGone)
			fmt.Fprintf(handler.output, "process '%v': wait error: process killed\n", pid)
			return
		}
		if errors.Is(err, models.ErrProcessWaitTimeout) {
			http.Error(w, err.Error(), http.StatusRequestTimeout)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(handler.output, "process '%v': wait error: unexpected error\n", pid)
		return
	}
	if err := httpx.WriteJson(w, http.StatusOK, command); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(handler.output, "process '%v': wait error: unexpected error\n", pid)
		return
	}
	fmt.Fprintf(handler.output, "process '%v': received command '%v'\n", pid, command)
}

func (handler *processesBrokerHttpHandler) list(w http.ResponseWriter, r *http.Request) {
	if err := httpx.AsNothing(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	processes, err := handler.processesBroker.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := httpx.WriteJson(w, http.StatusOK, models.MakeProcessesBody(processes...)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(handler.output, "list: unexpected error\n")
		return
	}
	fmt.Fprintf(handler.output, "list: success\n")
}

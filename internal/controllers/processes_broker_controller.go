package controllers

import (
	"context"
	"sync"
	"time"

	"github.com/mainden/stdhttp/internal/models"
)

type processesBrokerController struct {
	waitTimeout     time.Duration
	processes       map[int]models.ProcessModel
	processesBus    map[int]chan string
	processesExpire map[int]time.Time
	processesMutex  *sync.RWMutex
}

func NewProcessesBrokerController(waitTimeout time.Duration) *processesBrokerController {
	if waitTimeout <= time.Second {
		waitTimeout = 10 * time.Second
	}
	return &processesBrokerController{
		waitTimeout:     waitTimeout,
		processes:       make(map[int]models.ProcessModel),
		processesBus:    make(map[int]chan string),
		processesExpire: make(map[int]time.Time),
		processesMutex:  &sync.RWMutex{},
	}
}

func (controller *processesBrokerController) extendExpire(pid int) {
	controller.processesExpire[pid] = time.Now().Add(controller.waitTimeout + time.Second)
}

func (controller *processesBrokerController) extendExpireW(pid int) {
	controller.processesMutex.Lock()
	defer controller.processesMutex.Unlock()
	controller.extendExpire(pid)
}

func (controller *processesBrokerController) exists(pid int) bool {
	_, ok := controller.processes[pid]
	return ok
}

func (controller *processesBrokerController) existsR(pid int) bool {
	controller.processesMutex.RLock()
	defer controller.processesMutex.RUnlock()
	return controller.exists(pid)
}

func (controller *processesBrokerController) Register(ctx context.Context, process models.ProcessModel) error {
	if controller.existsR(process.Pid) {
		return models.ErrProcessExists
	}

	controller.processesMutex.Lock()
	defer controller.processesMutex.Unlock()
	if controller.exists(process.Pid) {
		return models.ErrProcessExists
	}

	controller.processes[process.Pid] = process
	controller.extendExpire(process.Pid)
	controller.processesBus[process.Pid] = make(chan string, 1)
	return nil
}

func (controller *processesBrokerController) bus(pid int) (chan string, bool) {
	bus, ok := controller.processesBus[pid]
	return bus, ok
}

func (controller *processesBrokerController) busR(pid int) (chan string, bool) {
	controller.processesMutex.RLock()
	defer controller.processesMutex.RUnlock()
	return controller.bus(pid)
}

func (controller *processesBrokerController) kill(pid int) error {
	if !controller.exists(pid) {
		return models.ErrProcessNotFound
	}
	bus, ok := controller.bus(pid)
	if !ok {
		return models.ErrProcessNotFound
	}
	if bus == nil {
		return nil
	}
	close(bus)
	controller.processesBus[pid] = nil
	delete(controller.processes, pid)
	return nil
}

func (controller *processesBrokerController) Kill(ctx context.Context, pid int) error {
	if !controller.existsR(pid) {
		return models.ErrProcessNotFound
	}
	bus, ok := controller.busR(pid)
	if !ok {
		return models.ErrProcessNotFound
	}
	if bus == nil {
		return nil
	}

	controller.processesMutex.Lock()
	defer controller.processesMutex.Unlock()
	return controller.kill(pid)
}

func (controller *processesBrokerController) SendCommand(ctx context.Context, pid int, command string) error {
	bus, ok := controller.busR(pid)
	if !ok || bus == nil {
		return models.ErrProcessNotFound
	}

	controller.processesMutex.Lock()
	defer controller.processesMutex.Unlock()
	bus, ok = controller.bus(pid)
	if !ok || bus == nil {
		return models.ErrProcessNotFound
	}
	select {
	case bus <- command:
		return nil
	default:
		return models.ErrProcessBusy
	}
}

func (controller *processesBrokerController) WaitCommand(ctx context.Context, pid int) (string, error) {
	bus, ok := controller.busR(pid)
	if !ok {
		return "", models.ErrProcessNotFound
	}
	if bus == nil {
		return "", models.ErrProcessKilled
	}
	controller.extendExpireW(pid)
	ctx, cancel := context.WithTimeout(ctx, controller.waitTimeout)
	defer cancel()
	select {
	case message, ok := <-bus:
		if ok {
			return message, nil
		}
		return "", models.ErrProcessKilled
	case <-ctx.Done():
		return "", models.ErrProcessWaitTimeout
	}
}

func (controller *processesBrokerController) List(ctx context.Context) (models.ProcessModels, error) {
	controller.processesMutex.RLock()
	defer controller.processesMutex.RUnlock()
	processes := make(models.ProcessModels, 0, len(controller.processes))
	for _, process := range controller.processes {
		process.Expired = controller.processesExpire[process.Pid].Before(time.Now())
		processes = append(processes, process)
	}
	return processes, nil
}

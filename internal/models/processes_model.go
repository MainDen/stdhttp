package models

import (
	"errors"
)

var (
	ErrProcessExists      = errors.New("process exists")
	ErrProcessNotFound    = errors.New("process not found")
	ErrProcessKilled      = errors.New("process killed")
	ErrProcessWaitTimeout = errors.New("process wait timeout")
	ErrProcessBusy        = errors.New("process busy")
)

type ProcessModel struct {
	Pid         int
	ClientName  string
	CommandName string
	CommandArgs []string
	Persistent  bool
	Expired     bool
}

type ProcessModels []ProcessModel

//go:build !windows

package execx

import "os/exec"

type ProcessGroup struct{}

func NewProcessGroup() *ProcessGroup {
	return &ProcessGroup{}
}

func (pg *ProcessGroup) Add(cmd *exec.Cmd) error {
	// Handled by SIGTERM
	return nil
}

func (pg *ProcessGroup) Close() error {
	// Handled by SIGTERM
	return nil
}

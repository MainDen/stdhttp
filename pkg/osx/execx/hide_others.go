//go:build !windows

package execx

import (
	"os/exec"
)

func CmdHide(cmd *exec.Cmd) {}

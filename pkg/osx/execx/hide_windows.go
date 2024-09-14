package execx

import (
	"os/exec"
	"syscall"
)

func CmdHide(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
}

package execx

import (
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procCreateJobObject          = kernel32.NewProc("CreateJobObjectW")
	procSetInformationJobObject  = kernel32.NewProc("SetInformationJobObject")
	procAssignProcessToJobObject = kernel32.NewProc("AssignProcessToJobObject")
)

type JOBOBJECT_EXTENDED_LIMIT_INFORMATION struct {
	BasicLimitInformation struct {
		PerProcessUserTimeLimit uint64
		PerJobUserTimeLimit     uint64
		LimitFlags              uint32
		MinimumWorkingSetSize   uintptr
		MaximumWorkingSetSize   uintptr
		ActiveProcessLimit      uint32
		Affinity                uintptr
		PriorityClass           uint32
		SchedulingClass         uint32
	}
	IoInfo struct {
		ReadOperationCount  uint64
		WriteOperationCount uint64
		OtherOperationCount uint64
		ReadTransferCount   uint64
		WriteTransferCount  uint64
		OtherTransferCount  uint64
	}
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

const (
	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE = 0x2000

	PROCESS_ALL_ACCESS = 0x1F0FFF
)

func CreateJobObject() (syscall.Handle, error) {
	handle, _, err := procCreateJobObject.Call(0, 0)
	if handle == 0 {
		return 0, err
	}
	return syscall.Handle(handle), nil
}

func SetJobObjectLimit(job syscall.Handle) error {
	var info JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	info.BasicLimitInformation.LimitFlags = JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE

	ret, _, err := procSetInformationJobObject.Call(uintptr(job), 9, uintptr(unsafe.Pointer(&info)), uintptr(unsafe.Sizeof(info)))
	if ret == 0 {
		return err
	}
	return nil
}

func AssignProcessToJobObject(job syscall.Handle, process syscall.Handle) error {
	ret, _, err := procAssignProcessToJobObject.Call(uintptr(job), uintptr(process))
	if ret == 0 {
		return err
	}
	return nil
}

type ProcessGroup struct {
	job syscall.Handle
}

func NewProcessGroup() *ProcessGroup {
	return &ProcessGroup{}
}

func (pg *ProcessGroup) Add(cmd *exec.Cmd) error {
	if pg.job == 0 {
		job, err := CreateJobObject()
		if err != nil {
			return fmt.Errorf("failed to create job object: %w", err)
		}
		if err := SetJobObjectLimit(job); err != nil {
			return fmt.Errorf("failed to set job object limit: %w", err)
		}
		pg.job = job
	}
	handle, err := syscall.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(cmd.Process.Pid))
	if err != nil {
		return fmt.Errorf("failed to open process: %w", err)
	}
	if err := AssignProcessToJobObject(pg.job, handle); err != nil {
		return fmt.Errorf("failed to assign process to job object: %w", err)
	}
	return nil
}

func (pg *ProcessGroup) Close() error {
	if pg.job != 0 {
		return syscall.CloseHandle(pg.job)
	}
	return nil
}

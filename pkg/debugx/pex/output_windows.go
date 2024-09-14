//go:build windows

package pex

import (
	"io"
	"syscall"
	"unsafe"
)

const (
	MB_OK              = 0x00000000
	MB_ICONERROR       = 0x00000010
	MB_ICONINFORMATION = 0x00000040
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procMessageBoxW = user32.NewProc("MessageBoxW")
)

func IconError() int {
	return MB_ICONERROR
}

func IconInformation() int {
	return MB_ICONINFORMATION
}

func MessageBox(icon int, title string, text string) error {
	t, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return err
	}
	msg, err := syscall.UTF16PtrFromString(text)
	if err != nil {
		return err
	}

	_, _, err = procMessageBoxW.Call(
		0,                            // HWND (handle to the window)
		uintptr(unsafe.Pointer(msg)), // Message text
		uintptr(unsafe.Pointer(t)),   // Title text
		uintptr(MB_OK|icon),          // Style of the message box
	)
	return err
}

type guiOutput struct {
	icon  int
	title string
}

func (o *guiOutput) Write(p []byte) (int, error) {
	if err := MessageBox(o.icon, o.title, string(p)); err != nil {
		return 0, err
	}
	return len(p), nil
}

func GUIOutput(icon int, title string) io.Writer {
	return &guiOutput{
		icon:  icon,
		title: title,
	}
}

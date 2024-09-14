//go:build windows

package pex

import (
	"debug/pe"
	"os"
)

func isGUI() bool {
	fileName, err := os.Executable()
	if err != nil {
		return false
	}
	f, err := pe.Open(fileName)
	if err != nil {
		return false
	}
	defer f.Close()
	var subsystem uint16
	if header, ok := f.OptionalHeader.(*pe.OptionalHeader64); ok {
		subsystem = header.Subsystem
	} else if header, ok := f.OptionalHeader.(*pe.OptionalHeader32); ok {
		subsystem = header.Subsystem
	}
	return subsystem == pe.IMAGE_SUBSYSTEM_WINDOWS_GUI
}

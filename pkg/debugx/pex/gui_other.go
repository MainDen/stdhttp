//go:build !windows

package pex

func isGUI() bool {
	return false
}

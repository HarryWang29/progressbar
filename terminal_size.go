//go:build !windows
// +build !windows

package progressbar

import (
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

type windowSize struct {
	rows uint16
	cols uint16
}

var out *os.File
var err error
var sz windowSize

func getTermSize(fd int) (int, int) {
	width, height, _ := terminal.GetSize(fd)
	return width, height
}

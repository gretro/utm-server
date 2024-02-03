//go:build !darwin
// +build !darwin

package detectos

import (
	"os"

	"github.com/gretro/utm_server/src/system"
)

func AssertDarwin() {
	l := system.GetComponentLogger("detectos")
	l.Error("This executable can only be run on macOS systems.")

	os.Exit(1)
}

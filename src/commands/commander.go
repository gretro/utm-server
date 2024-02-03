package commands

import (
	"errors"
	"log/slog"
	"os/exec"
	"path"
	"strings"

	"github.com/gretro/utm_server/src/config"
	"github.com/gretro/utm_server/src/system"
)

const (
	execSuffix = "Contents/MacOS/utmctl"

	machineNotFoundErrorText = "Error: Virtual machine not found."
)

var (
	ErrMachineNotFound = errors.New("machine not found")
	ErrCommandFailed   = errors.New("command failed")
)

type Commander struct {
	l       *slog.Logger
	cmdPath string
}

func New(appConfig *config.AppConfig) *Commander {
	return &Commander{
		l:       system.GetComponentLogger("commander"),
		cmdPath: path.Join(appConfig.UTMPath, execSuffix),
	}
}

func (c *Commander) exec(args []string) (string, error) {
	cmd := exec.Command(c.cmdPath, args...)
	c.l.Debug("Executing command", "cmd", cmd.String())

	out, err := cmd.CombinedOutput()
	if err != nil {
		c.l.Error("Failed to execute command", "cmd", cmd.String(), system.ErrorLabel, err)
		return "", wrapCmdError(c.l, err)
	}

	c.l.Debug("Command executed", "cmd", cmd.String(), "exitCode", cmd.ProcessState.ExitCode())
	return strings.TrimSpace(string(out)), nil
}

func wrapCmdError(l *slog.Logger, err error) error {
	exitErr := &exec.ExitError{}

	if errors.As(err, &exitErr) {
		stdErr := strings.TrimSpace(string(exitErr.Stderr))

		l.Debug("Command std err", "stdErr", stdErr)
		if stdErr == machineNotFoundErrorText {
			return ErrMachineNotFound
		}
	}

	return ErrCommandFailed
}

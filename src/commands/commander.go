package commands

import (
	"errors"
	"fmt"
	"io"
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
	errorFromEventText       = "Error from event"
)

var (
	ErrMachineNotFound  = errors.New("machine not found")
	ErrCommandFailed    = errors.New("command failed")
	ErrInvalidOperation = errors.New("invalid operation")
)

type Commander struct {
	l       *slog.Logger
	cmdPath string
	vmPath  string
}

func New(appConfig *config.AppConfig) *Commander {
	return &Commander{
		l:       system.GetComponentLogger("commander"),
		cmdPath: path.Join(appConfig.UTMPath, execSuffix),
		vmPath:  appConfig.VMPath,
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
	strOut := strings.TrimSpace(string(out))

	if strings.HasPrefix(strOut, errorFromEventText) {
		c.l.Error("Error from event", "cmd", cmd.String(), "output", strOut)
		return "", ErrInvalidOperation
	}

	return strOut, nil
}

func wrapCmdError(l *slog.Logger, err error) error {
	exitErr := &exec.ExitError{}

	if errors.As(err, &exitErr) {
		stdErr := strings.TrimSpace(string(exitErr.Stderr))

		l.Error("Command failed", "exitcode", exitErr.ExitCode(), "stdErr", stdErr)
		if stdErr == machineNotFoundErrorText {
			return ErrMachineNotFound
		}
	}

	return ErrCommandFailed
}

func (c *Commander) pipedExec(cmds ...*exec.Cmd) ([]byte, error) {
	var prevOutPipe io.ReadCloser
	var err error

	// Chain all commands together
	for i, cmd := range cmds {
		if prevOutPipe != nil {
			cmd.Stdin = prevOutPipe
		}

		prevOutPipe, err = cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("failed to create out pipe for command %d: %w", i, err)
		}
	}

	// Starting all commands
	for _, cmd := range cmds {
		err = cmd.Start()
		if err != nil {
			return nil, wrapCmdError(c.l, err)
		}
	}

	// Listening for the last output in a separate goroutine
	var out []byte
	done := make(chan error)
	go func() {
		out, err = io.ReadAll(prevOutPipe)
		done <- err
	}()

	// Waiting for all commands to finish
	for _, cmd := range cmds {
		err = cmd.Wait()
		if err != nil {
			return nil, wrapCmdError(c.l, err)
		}
	}

	err = <-done
	if err != nil {
		return nil, err
	}

	return out, nil
}

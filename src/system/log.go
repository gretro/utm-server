package system

import "log/slog"

const (
	ErrorLabel = "err"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.Default().Handler())
}

func GetComponentLogger(component string) *slog.Logger {
	return logger.With("component", component)
}

func SystemLogger() *slog.Logger {
	return logger
}

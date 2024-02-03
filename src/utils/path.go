package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ResolvePath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}

		path = filepath.Join(homeDir, path[2:])
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	return path, nil
}

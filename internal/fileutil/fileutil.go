package fileutil

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func ExpandPath(path string) (string, error) {
	if path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("expanding path %w", err)
		}

		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

// CheckFileExists checks if a file (not a directory) can be found at the given path.
// It returns false for directories even if they exist.
func CheckFileExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fileInfo.IsDir() {
		return false
	}

	return true
}

// EnsureDirectory creates the directory if it doesn't exist or verifies it's a directory
func EnsureDirectory(path string, perm os.FileMode) error {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// Directory doesn't exist, create it
			if err := os.MkdirAll(path, perm); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to check directory: %w", err)
	}

	// Path exists but isn't a directory
	if !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}

	return nil
}

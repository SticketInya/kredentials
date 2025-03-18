package fileutil

import (
	"fmt"
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

// Checks if a file (not a directory) can be found at the given path
// Returns 'false' for directories even if they exists.
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

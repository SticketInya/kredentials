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

package helpers

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func AtProjectRoot(path string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", path)
}

// ExpandPath takes a path and if it has a "~/" prefix, it will expand
// it to os.UserHomeDir()
func ExpandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, path[2:]), nil
}

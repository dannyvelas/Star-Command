package helpers

import (
	"path/filepath"
	"runtime"
)

func AtProjectRoot(path string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", path)
}

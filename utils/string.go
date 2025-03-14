package utils

import (
	"path/filepath"
	"runtime"
	"strings"
)

func TrimSuffixSlashOnPaths(paths ...string) string {
	concatPath := strings.ReplaceAll(strings.Join(paths, "/"), "//", "/")
	return strings.TrimSuffix(concatPath, "/")
}

func ConcatPaths(v ...string) (values string) {
	if len(v) > 0 {
		values = v[0]
	}

	separator := string(filepath.Separator)

	if runtime.GOOS == "windows" {
		separator = "\\"
	}

	for _, i := range v[1:] {
		if !strings.HasSuffix(values, separator) {
			values += separator
		}
		values += i
	}

	return values
}

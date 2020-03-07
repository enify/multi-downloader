package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IsFileExist check if file exists
func IsFileExist(path string) bool {
	filePath, _ := filepath.Abs(path)
	_, err := os.Stat(filePath)

	return !os.IsNotExist(err)
}

// GenUnusedFilename generate a unused filename in path
func GenUnusedFilename(path, filename string) string {
	suffix := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, suffix)

	for i := 1; i <= 100; i++ {
		newName := fmt.Sprintf("%s(%d)%s", name, i, suffix)
		if !IsFileExist(filepath.Join(path, newName)) {
			return newName
		}
	}
	return filename
}

// FollowPathRule let string fit OS path rule
func FollowPathRule(s string) (path string) {
	exp := regexp.MustCompile(`[/\\:*?"<>|]+`)
	path = exp.ReplaceAllString(s, "_")

	return
}

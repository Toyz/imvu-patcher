package core

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func toInt(s string) int {
	i, _ := strconv.Atoi(s)

	return i
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func noExtension(filename string) string {
	var extension = filepath.Ext(filename)
	return strings.TrimRight(filename, extension)
}

package utils

import (
	"path/filepath"
	"strings"
)

const TemplateDirName = "templates"

func GetTemplNameByFile(filename string) string {
	ext := filepath.Ext(filename)
	return strings.Replace(filename, ext, "", 1)
}

package templates

import (
	"path/filepath"
	"strings"
)

const DirName = "templates"

func GetTemplNameByFile(filename string) string {
	ext := filepath.Ext(filename)
	return strings.Replace(filename, ext, "", 1)
}

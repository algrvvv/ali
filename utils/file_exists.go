package utils

import (
	"errors"
	"os"

	"github.com/algrvvv/ali/logger"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)

	// файл существует
	if err == nil {
		return true
	}

	// файла не существует
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	logger.SaveDebugf("failed to get file stat: %v", err)
	return false
}

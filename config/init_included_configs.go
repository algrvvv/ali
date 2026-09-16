package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/algrvvv/ali/v2/logger"
	"github.com/spf13/viper"
)

func InitIncluded() {
	include := viper.GetStringSlice("include")
	logger.SaveDebugf("includes: %v", include)

	for _, includePath := range include {
		if strings.HasPrefix(includePath, "~") {
			logger.SaveDebugf("user use ~ in dir param")
			home, err := os.UserHomeDir()
			if err != nil {
				logger.SaveDebugf("failed to get user home dir")
				return
			}

			logger.SaveDebugf("got home user dir: %q", home)
			includePath = strings.Replace(includePath, "~", home, 1)
			logger.SaveDebugf("command dir after change ~ to home dir: %q", includePath)
		}
		logger.SaveDebugf("include path: %s", includePath)

		info, err := os.Stat(includePath)
		if err != nil {
			fmt.Println("failed to get include file stat: ", err)
			logger.SaveDebugf("failed to get include file stat: %v", err)
			return
		}

		path := includePath
		if info.IsDir() {
			path = filepath.Join(path, ".ali")
		}

		logger.SaveDebugf("final unclude path: %s", path)
		viper.SetConfigFile(path)
		if err := viper.MergeInConfig(); err != nil {
			logger.SaveDebugf("uncluded config not found")
		} else {
			logger.SaveDebugf("uncluded config loaded")
		}
	}
}

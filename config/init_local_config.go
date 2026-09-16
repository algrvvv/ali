package config

import (
	"github.com/algrvvv/ali/v2/logger"
	"github.com/algrvvv/ali/v2/utils"
	"github.com/spf13/viper"
)

const localConfig = ".ali"

func InitLocal(lviper *viper.Viper, localEnv bool) {
	logger.SaveDebugf("local config: %s", localConfig)

	lviper = viper.New()
	lviper.SetConfigName(localConfig)
	lviper.SetConfigType(utils.YamlConfigurationType)
	lviper.AddConfigPath(".")

	viper.SetConfigName(localConfig)
	viper.SetConfigType(utils.YamlConfigurationType)
	viper.AddConfigPath(".")

	if err := lviper.ReadInConfig(); err != nil {
		logger.SaveDebugf("load local config error: %v", err)
		// utils.CheckError(err)
	} else {
		logger.SaveDebugf("local viper read config successfully")
	}

	if localEnv {
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			logger.SaveDebugf("load local config error: %v", err)
			// utils.CheckError(err)
		}
		logger.SaveDebugf("using config: %s", viper.ConfigFileUsed())
	} else {
		if err := viper.MergeInConfig(); err != nil {
			logger.SaveDebugf("local config not found")
		} else {
			logger.SaveDebugf("local config loaded")
		}
	}
}

package utils

import (
	"maps"

	"github.com/algrvvv/ali/logger"
	"github.com/spf13/viper"
)

func GetEnvs(alias *AliasEntry) map[string]any {
	// здесь мы получаем и глобальные переменные
	// окружения и для конкретной команды (алиаса)
	logger.SaveDebugf("search envs for %s", alias.AliasName)

	env := viper.GetStringMap("env")
	maps.Copy(env, alias.Env)

	return env
}

package envs

import (
	"maps"

	"github.com/algrvvv/ali/v2/aliases"
	"github.com/algrvvv/ali/v2/logger"
	"github.com/spf13/viper"
)

func GetAliasEnvs(alias *aliases.AliasEntry) map[string]any {
	// здесь мы получаем и глобальные переменные
	env := viper.GetStringMap("env")

	if alias != nil {
		// окружения и для конкретной команды (алиаса)
		logger.SaveDebugf("search envs for %s", alias.AliasName)
		maps.Copy(env, alias.Env)
	}

	return env
}

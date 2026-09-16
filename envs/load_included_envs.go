package envs

import (
	"github.com/algrvvv/ali/v2/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// loadIncluded загрузка внешних енвов
func loadIncluded() map[string]any {
	envFiles := viper.GetStringSlice("envs")
	logger.SaveDebugf("found included env files: %v", envFiles)

	includedEnvs := make(map[string]any)

	// просто вариант: godotenv.Load(envFiles...)
	// но нам не подходит, так как будет выброс исключения при ошибке загрузке енва
	// в нашем случае, необходимо тихо это обработать и продолжить работу
	for _, envFile := range envFiles {
		envMap, err := godotenv.Read()

		if err != nil {
			// тихо обрабатываем ошибку
			logger.SaveDebugf("Error loading %q file: %v", envFile, err)
		}

		// делаем каст string -> any
		for k, v := range envMap {
			includedEnvs[k] = v
		}
	}

	return includedEnvs
}

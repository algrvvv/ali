package utils

import "github.com/spf13/viper"

func GetAlias(aliasName string) string {
	key := "aliases." + aliasName

	return viper.GetString(key)
}

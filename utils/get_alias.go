package utils

import "github.com/spf13/viper"

func GetAlias(aliasName string) any {
	key := "aliases." + aliasName

	return viper.Get(key)
}

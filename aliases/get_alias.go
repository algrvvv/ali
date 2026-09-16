package aliases

import "github.com/spf13/viper"

func GetByName(aliasName string) any {
	key := "aliases." + aliasName

	return viper.Get(key)
}

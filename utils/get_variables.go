package utils

import (
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

type varConfig struct {
	Vars map[string]string `mapstructure:"vars"`
}

func getVars() (map[string]string, error) {
	var vars varConfig

	err := viper.Unmarshal(&vars)
	if err != nil {
		return nil, err
	}

	return vars.Vars, nil
}

func GetVariables(input string, vars map[string]string) string {
	re := regexp.MustCompile(`\{\{(\w+)\}\}`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		// key := re.FindStringSubmatch(match)[1]
		key := strings.ToLower(re.FindStringSubmatch(match)[1])
		if val, ok := vars[key]; ok {
			return val
		}
		return match // если переменная не найдена — не трогать
	})
}

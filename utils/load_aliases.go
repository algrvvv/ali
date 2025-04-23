package utils

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type AliasEntry struct {
	AliasName string         `mapstructure:"alias"`
	Aliases   []string       `mapstructure:"aliases"`
	Cmds      []string       `mapstructure:"cmds"`
	Desc      string         `mapstructure:"desc"`
	Env       map[string]any `mapstructure:"env"`
	Parallel  bool           `mapstructure:"parallel"`
	// TODO: реализовать
	Dir string `mapstructure:"dir"`
}

func LoadAliases(v *viper.Viper) map[string]AliasEntry {
	raw := v.GetStringMap("aliases")
	out := make(map[string]AliasEntry)

	for key, val := range raw {
		switch v := val.(type) {
		case string:
			out[key] = AliasEntry{
				AliasName: key,
				Cmds:      []string{v},
			}
		case map[string]any:
			var entry AliasEntry
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				Result:  &entry,
				TagName: "mapstructure",
			})
			if err != nil {
				// WARN: не забыть добавить обработку ошибки
				panic(err)
			}

			err = decoder.Decode(v)
			if err != nil {
				// WARN: не забыть добавить обработку ошибки
				panic(err)
			}

			entry.AliasName = key
			out[key] = entry
		default:
			fmt.Printf("unsupported alias value type for %q: %T\n", key, v)
		}
	}

	return out
}

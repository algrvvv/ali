package envs

import (
	"errors"
	"maps"

	"github.com/algrvvv/ali/v2/aliases"
	"github.com/algrvvv/ali/v2/utils"
)

func Load(entry *aliases.AliasEntry, hiddenFilepath string) map[string]any {
	hiddenEnvs, err := LoadHiddenEnv(hiddenFilepath)
	if err != nil && !errors.Is(err, ErrHiddenEnvFileNotFound) {
		utils.PrintError("failed to load hidden envs", nil)
	}
	envs := GetEnvs(entry)
	maps.Copy(envs, hiddenEnvs)

	return envs
}

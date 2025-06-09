package utils

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrHiddenEnvFileNotFound = errors.New("hidden envs file not found")

func LoadHiddenEnv(filepath string) (map[string]any, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrHiddenEnvFileNotFound
		}
		return nil, err
	}

	var res map[string]any
	err = yaml.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

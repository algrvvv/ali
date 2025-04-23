package utils

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/algrvvv/ali/logger"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

const (
	YamlConfigurationType = "yaml"
	JsonConfigurationType = "json"
	TomlConfigurationType = "toml"
)

var ErrUnsupportedConfigType = errors.New("unsupported config type")

// GetConfigurationType функция для определения типа конфигурации.
// Deprecated: функция больше не используется
func GetConfigurationType(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TomlConfigurationType, err
	}

	if json.Valid(data) {
		return JsonConfigurationType, nil
	}
	logger.SaveDebugf("failed to read config [%s] as json", path)

	var yamlData map[string]any
	if err = yaml.Unmarshal(data, &yamlData); err == nil {
		return YamlConfigurationType, nil
	}
	logger.SaveDebugf("failed to read config [%s] as yaml: %v", path, err)

	var tomlData map[string]any
	if err = toml.Unmarshal(data, &tomlData); err == nil {
		return TomlConfigurationType, nil
	}
	logger.SaveDebugf("failed to read config [%s] as toml: %v", path, err)

	return "", ErrUnsupportedConfigType
}

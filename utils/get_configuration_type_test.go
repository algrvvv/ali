package utils_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/algrvvv/ali/utils"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

func TestGetConfigurationTOMLType(t *testing.T) {
	t.Parallel()

	tomlData, err := toml.Marshal(map[string]any{
		"aliases": map[string]any{
			"test": "echo 'hello world'",
		},
	})
	if err != nil {
		t.Errorf("failed to prepare toml data: %s", err)
	}
	t.Log("data for toml confgiration type file created")

	tomlTmp, err := os.CreateTemp("", "get-configration-type-test-*.toml")
	if err != nil {
		t.Errorf("failed to create temp file for toml data: %v", err)
	}
	t.Log("toml temp confgiration file created")

	_, err = tomlTmp.Write(tomlData)
	if err != nil {
		t.Errorf("failed to write toml data to temp toml data file: %v", err)
	}
	t.Log("data writed to toml temp file")
	defer os.Remove(tomlTmp.Name())

	wantTomlType, err := utils.GetConfigurationType(tomlTmp.Name())
	if err != nil {
		t.Errorf("failed to get first confgiration type: %v", err)
	}

	if wantTomlType == utils.TomlConfigurationType {
		t.Logf("SUCCESS! want: %s; get type: %s", utils.TomlConfigurationType, wantTomlType)
	} else {
		t.Errorf("FAIL! want: %s; get type: %s", utils.TomlConfigurationType, wantTomlType)
	}
}

func TestGetConfigurationYAMLType(t *testing.T) {
	t.Parallel()

	yamlData, err := yaml.Marshal(map[string]any{
		"aliases": map[string]any{
			"test": "echo 'hello world'",
		},
	})
	if err != nil {
		t.Errorf("failed to prepare yaml data: %s", err)
	}
	t.Log("data for yaml confgiration type file created")

	yamlTmp, err := os.CreateTemp("", "get-configration-type-test-*.yml")
	if err != nil {
		t.Errorf("failed to create temp file for yaml data: %v", err)
	}
	t.Log("yaml temp confgiration file created")

	_, err = yamlTmp.Write(yamlData)
	if err != nil {
		t.Errorf("failed to write yaml data to temp yaml data file: %v", err)
	}
	t.Log("data writed to yaml temp file")
	defer os.Remove(yamlTmp.Name())

	wantYamlType, err := utils.GetConfigurationType(yamlTmp.Name())
	if err != nil {
		t.Errorf("failed to get first confgiration type: %v", err)
	}

	if wantYamlType == utils.YamlConfigurationType {
		t.Logf("SUCCESS! want: %s; get type: %s", utils.YamlConfigurationType, wantYamlType)
	} else {
		t.Errorf("FAIL! want: %s; get type: %s", utils.YamlConfigurationType, wantYamlType)
	}
}

func TestGetConfigurationJSONType(t *testing.T) {
	t.Parallel()

	jsonData, err := json.Marshal(map[string]any{
		"aliases": map[string]any{
			"test": "echo 'hello world'",
		},
	})
	if err != nil {
		t.Errorf("failed to prepare json data: %s", err)
	}
	t.Log("data for json confgiration type file created")

	jsonTmp, err := os.CreateTemp("", "get-configration-type-test-*.json")
	if err != nil {
		t.Errorf("failed to create temp file for json data: %v", err)
	}
	t.Log("json temp confgiration file created")

	_, err = jsonTmp.Write(jsonData)
	if err != nil {
		t.Errorf("failed to write json data to temp json data file: %v", err)
	}
	t.Log("data writed to json temp file")
	defer os.Remove(jsonTmp.Name())

	wantJsonType, err := utils.GetConfigurationType(jsonTmp.Name())
	if err != nil {
		t.Errorf("failed to get first confgiration type: %v", err)
	}

	if wantJsonType == utils.JsonConfigurationType {
		t.Logf("SUCCESS! want: %s; get type: %s", utils.JsonConfigurationType, wantJsonType)
	} else {
		t.Errorf("FAIL! want: %s; get type: %s", utils.JsonConfigurationType, wantJsonType)
	}
}

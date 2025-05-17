package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	PluginsDirName   = "plugins"
	PluginConfigName = "ali-plugin.yml"
)

func GetPluginsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".ali")

	pluginsDirPath := filepath.Join(configDir, PluginsDirName)
	return pluginsDirPath, nil
}

func GetPluginDirByName(name string) (string, error) {
	plugins, err := GetPluginsDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(plugins, name), nil
}

func GetPlugins() ([]string, error) {
	pluginDir, err := GetPluginsDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		return nil, err
	}

	var plugins []string
	for idx, entry := range entries {
		pluginPath := filepath.Join(pluginDir, entry.Name())
		pluginConfigPath := filepath.Join(pluginPath, PluginConfigName)

		var pluginDesc string
		var pluginStatus string

		v, err := GetViperForPlugin(pluginConfigPath)
		if err == nil {
			pluginDesc = fmt.Sprintf("\n%s%s%s", Colors["gray"], v.GetString("desc"), Colors["reset"])
			if v.GetString("exec") == "" {
				pluginStatus = fmt.Sprintf("%s[!] plugin disabled: invalid exec param%s", Colors["red"], Colors["reset"])
			}
		}

		if _, err := os.Stat(pluginConfigPath); err != nil {
			pluginStatus = fmt.Sprintf("%s[!] plugin disabled: plugin config not found%s", Colors["red"], Colors["reset"])
		}

		if pluginStatus == "" {
			pluginStatus = fmt.Sprintf("%s[+] ok%s", Colors["green"], Colors["reset"])
		}

		e := fmt.Sprintf("%d. %-20s(%s)\t%s%s",
			idx+1,
			entry.Name(),
			pluginPath,
			pluginStatus,
			pluginDesc,
		)

		plugins = append(plugins, e)
	}

	return plugins, nil
}

func GetViperForPlugin(path string) (*viper.Viper, error) {
	pluginViper := viper.New()
	pluginViper.SetConfigFile(path)
	if err := pluginViper.ReadInConfig(); err != nil {
		return nil, err
	}
	return pluginViper, nil
}

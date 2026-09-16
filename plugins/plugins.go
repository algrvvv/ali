package plugins

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/algrvvv/ali/v2/utils"
	"github.com/spf13/viper"
)

const (
	DirName    = "plugins"
	ConfigName = "ali-plugin.yml"
)

func GetPluginsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".ali")

	pluginsDirPath := filepath.Join(configDir, DirName)
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
		pluginConfigPath := filepath.Join(pluginPath, ConfigName)

		var pluginDesc string
		var pluginStatus string

		v, err := GetViperForPlugin(pluginConfigPath)
		if err == nil {
			pluginDesc = fmt.Sprintf("\n%s%s%s", utils.Colors["gray"], v.GetString("desc"), utils.Colors["reset"])
			if v.GetString("exec") == "" {
				pluginStatus = fmt.Sprintf("%s[!] plugin disabled: invalid exec param%s", utils.Colors["red"], utils.Colors["reset"])
			}
		}

		if _, err := os.Stat(pluginConfigPath); err != nil {
			pluginStatus = fmt.Sprintf("%s[!] plugin disabled: plugin config not found%s", utils.Colors["red"], utils.Colors["reset"])
		}

		if pluginStatus == "" {
			pluginStatus = fmt.Sprintf("%s[+] ok%s", utils.Colors["green"], utils.Colors["reset"])
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

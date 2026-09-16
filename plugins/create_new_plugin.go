package plugins

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/algrvvv/ali/v2/utils"
	"github.com/spf13/viper"
)

func CreateNew(name string) {
	// NOTE: плагин сам по себе будет папкой, которая будет содержать
	// конфигурацию плагина - файл ali-plugin.yml
	// и сами файлы плагина (на усмотрение каждого)
	//
	// NOTE: ali-plugin.yml будет иметь в себе:
	//  - desc: описание, которое будет по желание пользователя (no desc - по умолчанию)
	//  - exec: то, как запускать плагин (дефолтная директория, та в которой и будет лежать плагин)
	// название плагина будет равняться названию директории, в которой лежит плагин и его конфигурация
	//
	// WARN: если в директории плагина не будет конфигурации, то такой плагин будет невалидным

	fmt.Println("create new plugin: " + name)

	dir, err := GetPluginsDir()
	if err != nil {
		utils.PrintError("failed to get plugins dir", err)
		return
	}

	pluginDir := filepath.Join(dir, name)
	if _, err := os.Stat(pluginDir); err == nil {
		fmt.Printf("plugin with name %q already exists\n", name)
		return
	}

	err = os.MkdirAll(pluginDir, 0777)
	if err != nil {
		utils.PrintError("failed to create dir for plugin", err)
		return
	}

	// создаем пустой файл конфигурации
	pluginConfigPath := filepath.Join(pluginDir, ConfigName)
	f, err := os.Create(pluginConfigPath)
	if err != nil {
		utils.PrintError("failed to create plugin config", err)
		return
	}
	f.Close()

	// в файл конфигурации добавляем дефолтные занчения
	pluginConfigViper := viper.New()
	pluginConfigViper.SetConfigFile(pluginConfigPath)

	pluginConfigViper.Set("desc", "empty plugin description")
	pluginConfigViper.Set("exec", "echo 'command for start your plugin (exec from your plug dir)'")

	err = pluginConfigViper.WriteConfig()
	if err != nil {
		utils.PrintError("failed to write plugin config", err)
		return
	}

	fmt.Println("plugin location: ", pluginDir)
}

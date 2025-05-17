/*
Copyright © 2025 algrvvv <alexandrgr25@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/algrvvv/ali/logger"
	"github.com/algrvvv/ali/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// plugCmd represents the plug command
var (
	createNewPlug bool
	showPlugList  bool
	plugCmd       = &cobra.Command{
		Use:   "plug plugin-name",
		Short: "work with plugins",
		Run: func(cmd *cobra.Command, args []string) {
			if showPlugList {
				showPlugins()
				return
			}

			if len(args) < 1 {
				fmt.Println("failed to start `plug` command: expected plug name\nuse: ali plug pluginName --new; or ali plug --list")
				return
			}

			plugName := args[0]
			logger.SaveDebugf("got plugin name: %s", plugName)

			if createNewPlug {
				createNewPlugin(plugName)
				return
			}

			params := args[1:]
			execPlugin(plugName, params)
		},
	}
)

func showPlugins() {
	plugs, err := utils.GetPlugins()
	if err != nil {
		utils.PrintError("failed to get plugins list", err)
		return
	}

	if len(plugs) == 0 {
		fmt.Println("plugins not found")
		return
	}

	fmt.Println("plugins list:")
	for _, plug := range plugs {
		fmt.Println(plug)
	}
}

func execPlugin(name string, params []string) {
	logger.SaveDebugf("exec plugin: %s", name)

	pluginPath, err := utils.GetPluginDirByName(name)
	if err != nil {
		utils.PrintError("failed to get plugin", err)
		return
	}

	if _, err = os.Stat(pluginPath); err != nil {
		utils.PrintError("plugin not found", err)
		return
	}

	pluginConfigPath := filepath.Join(pluginPath, utils.PluginConfigName)
	if _, err = os.Stat(pluginConfigPath); err != nil {
		utils.PrintError("plugin disabled", err)
		return
	}

	v, err := utils.GetViperForPlugin(pluginConfigPath)
	if err != nil {
		utils.PrintError("failed to read plugin config", err)
		return
	}

	execCommand := v.GetString("exec")
	if execCommand == "" {
		fmt.Println("empty exec param for plugin")
		return
	}

	for _, param := range params {
		paramParts := strings.Split(param, "=")
		if len(paramParts) == 2 {
			execCommand += fmt.Sprintf(" -%s=%q", paramParts[0], paramParts[1])
		} else {
			execCommand += fmt.Sprintf(" -%s", param)
		}
	}
	logger.SaveDebugf("result command for exec: %s", execCommand)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe", "/C", execCommand)
	case "linux", "darwin":
		cmd = exec.Command("sh", "-c", execCommand)
	default:
		fmt.Println("unsuppored os")
		return
	}

	cmd.Dir = pluginPath
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		utils.PrintError("failed to start exec plugin command: %s", err)
		return
	}

	if err = cmd.Wait(); err != nil {
		utils.PrintError("failed to wait plugin command: %s", err)
		return
	}
}

func createNewPlugin(name string) {
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

	dir, err := utils.GetPluginsDir()
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
	pluginConfigPath := filepath.Join(pluginDir, utils.PluginConfigName)
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

func init() {
	rootCmd.AddCommand(plugCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// plugCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	plugCmd.Flags().BoolVar(&createNewPlug, "new", false, "create new plugin")
	plugCmd.Flags().BoolVar(&showPlugList, "list", false, "show list of plugins")
}

/*
Copyright © 2024 algrvvv <alexandrgr25@gmail.com>

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
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/algrvvv/ali/logger"
	"github.com/algrvvv/ali/utils"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup global config",
	Run: func(_ *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		utils.CheckError(err)
		configDir := filepath.Join(home, ".ali")
		err = os.MkdirAll(configDir, 0o700)
		if err != nil && !errors.Is(err, os.ErrExist) {
			utils.CheckError(err)
		}
		logger.SaveDebugf("got config dir: %s", configDir)

		configPath := filepath.Join(configDir, "config.toml")
		logger.SaveDebugf("got config path: %s", configPath)

		viper.SetConfigFile(configPath)
		viper.Set("aliases.test", "echo \"hello world\"")
		viper.Set("app.editor", "vi")
		viper.Set("app.default_config_type", utils.TomlConfigurationType)

		err = viper.WriteConfig()
		utils.CheckError(err)
		logger.SaveDebugf("config writed")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

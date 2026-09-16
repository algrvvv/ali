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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/algrvvv/ali/v2/logger"
	"github.com/algrvvv/ali/v2/utils"
)

const localConfig = ".ali"

// initCmd represents the init command
var (
	configFormat string = utils.YamlConfigurationType
	initCmd             = &cobra.Command{
		Use:   "init",
		Short: "Init new local config",
		Run: func(cmd *cobra.Command, args []string) {
			logger.SaveDebugf("save new local config with format as: %s", configFormat)

			dir, err := os.Getwd()
			utils.CheckError(err)
			logger.SaveDebugf("got dir: %s", dir)

			if utils.FileExists(localConfig) {
				fmt.Println("local config already exists")
				return
			}

			f, err := os.OpenFile(localConfig, os.O_CREATE|os.O_RDWR, 0o600)
			utils.CheckError(err)

			defer f.Close()
			fmt.Println("local config file initialized")

			switch configFormat {
			case utils.YamlConfigurationType:
				_, err = f.WriteString("aliases:\n")
			case utils.JsonConfigurationType:
				_, err = f.WriteString("{\n\t\"aliases\": {}\n}")
			case utils.TomlConfigurationType:
				_, err = f.WriteString("[aliases]\n")
			default:
				fmt.Println(utils.ErrUnsupportedConfigType)
				return
			}

			utils.CheckError(err)
			utils.CheckError(f.Close())
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)

	// NOTE: больше не поддерживается
	// initCmd.Flags().StringVarP(&configFormat, "format", "F", utils.YamlConfigurationType, "new local config type")
}

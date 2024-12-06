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
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/algrvvv/ali/logger"
	"github.com/algrvvv/ali/utils"
)

// editCmd represents the edit command
var (
	isLocal bool

	editCmd = &cobra.Command{
		Use:   "edit",
		Short: "A brief description of your command",
		Run: func(_ *cobra.Command, args []string) {
			editor := viper.GetString("app.editor")
			if editor == "" {
				editor = "vi"
			}
			logger.SaveDebugf("got editor: %s", editor)

			home, err := os.UserHomeDir()
			utils.CheckError(err)

			var path string
			if isLocal {
				var dir string
				dir, err = os.Getwd()
				utils.CheckError(err)

				path = filepath.Join(dir, ".ali")
			} else {
				path = filepath.Join(home, ".ali/config.toml")
			}
			logger.SaveDebugf("got path for edit config: %s", path)

			cmd := exec.Command(editor, path)
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			utils.CheckError(err)
		},
	}
)

func init() {
	rootCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	editCmd.Flags().BoolVarP(&isLocal, "local", "l", false, "edit local config file")
}

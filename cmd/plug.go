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

	"github.com/algrvvv/ali/v2/logger"
	"github.com/algrvvv/ali/v2/plugins"
	"github.com/spf13/cobra"
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
				plugins.Show()
				return
			}

			if len(args) < 1 {
				fmt.Println("failed to start `plug` command: expected plug name\nuse: ali plug pluginName --new; or ali plug --list")
				return
			}

			plugName := args[0]
			logger.SaveDebugf("got plugin name: %s", plugName)

			if createNewPlug {
				plugins.CreateNew(plugName)
				return
			}

			params := args[1:]
			plugins.Exec(plugName, params)
		},
	}
)

func init() {
	rootCmd.AddCommand(plugCmd)

	plugCmd.Flags().BoolVar(&createNewPlug, "new", false, "create new plugin")
	plugCmd.Flags().BoolVar(&showPlugList, "list", false, "show list of plugins")
}

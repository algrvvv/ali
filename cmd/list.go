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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/algrvvv/ali/utils"
)

const (
	color      = "\033[1;34m"
	resetColor = "\033[0m"
)

// listCmd represents the list command
var (
	fullPrint bool

	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "Get list aliases",
		Example: "ali list - to see all aliases\nali list push - to search aliases that contains 'push' in name",
		Run: func(_ *cobra.Command, args []string) {
			var search string
			if len(args) == 1 {
				search = args[0]
			}

			aliases, ok := viper.Get("aliases").(map[string]any)
			if !ok {
				fmt.Println("failed to get all aliases")
				return
			}

			fmt.Println("Available Aliases:")
			if fullPrint {
				count := len(aliases) - 1
				i := 0

				for alias, command := range aliases {
					if search == "" || (search != "" && strings.Contains(alias, search)) {
						prefix := "├──"
						if i == count {
							prefix = "└──"
						}

						fmt.Printf(
							"  %s %s%-6s%s -> %s\n",
							prefix,
							color,
							alias,
							resetColor,
							command,
						)
						i++
					}
				}

				return
			}

			fmt.Printf("+%s+%s+\n", strings.Repeat("-", 22), strings.Repeat("-", 42))
			fmt.Printf(
				"| Alias%s| Command%s|\n",
				strings.Repeat(" ", 22-len(" alias")),
				strings.Repeat(" ", 42-len(" command")),
			)
			fmt.Printf("+%s+%s+\n", strings.Repeat("-", 22), strings.Repeat("-", 42))

			for alias, command := range aliases {
				if search == "" || (search != "" && strings.Contains(alias, search)) {
					commandStr := command.(string)
					fmt.Printf(
						"| %s%-20s%s | %-40s |\n",
						color,
						utils.TruncateString(alias, 20),
						resetColor,
						utils.TruncateString(commandStr, 40),
					)
				}
			}

			fmt.Printf("+%s+%s+\n", strings.Repeat("-", 22), strings.Repeat("-", 42))
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().BoolVarP(&fullPrint, "full", "f", false, "use fast print")
}

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
	"errors"
	"fmt"

	"github.com/algrvvv/ali/v2/logger"
	"github.com/algrvvv/ali/v2/templates"
	"github.com/algrvvv/ali/v2/utils"
	"github.com/spf13/cobra"
)

var (
	createNewTempl bool
	forceAction    bool
	editTempl      bool
	showTemplList  bool
	templCmd       = &cobra.Command{
		Use:   "templ templName",
		Short: "use template configuration",
		Long:  `init new configuration by template`,
		Run: func(cmd *cobra.Command, args []string) {
			if showTemplList {
				if err := templates.Show(); err != nil {
					fmt.Println("failed to get list of templs")
					wrappedErr := fmt.Errorf("failed to show list of temps: %v", err)
					utils.CheckError(wrappedErr)
				}

				return
			}

			if len(args) != 1 {
				fmt.Println("failed to start `templ` command: expected templ name\nuse: ali templ templName --new; or ali templ --list")
				return
			}

			templName := args[0]
			logger.SaveDebugf("got templ name: %s", templName)

			if createNewTempl {
				if err := templates.CreateNew(templName, forceAction); err != nil {
					if errors.Is(err, templates.ErrTemplAlreadyExists) {
						fmt.Println(err.Error())
						return
					}

					fmt.Println(templates.ErrFailedToInitNewTempl.Error())
					// NOTE: оборачиываем ошибку для более удобного дальнейшего чтения
					// и либо просто говорим о том, что произошла ошибка, чекните логи,
					// либо выводим их
					wrappedErr := fmt.Errorf("failed to create new templ: %v", err)
					utils.CheckError(wrappedErr)
					return
				}
			}

			if editTempl {
				if err := templates.Edit(templName); err != nil {
					fmt.Println("failed to edit template")
					wrappedErr := fmt.Errorf("failed to edit templ: %v", err)
					utils.CheckError(wrappedErr)
				}

				return
			}

			if err := templates.Use(localViper, templName, forceAction); err != nil {
				fmt.Println("failed to use template")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(templCmd)

	templCmd.Flags().BoolVar(&createNewTempl, "new", false, "create new template")
	templCmd.Flags().BoolVar(&forceAction, "force", false, "create new template force (rewrite existed templ)")
	templCmd.Flags().BoolVar(&editTempl, "edit", false, "edit template by name")
	templCmd.Flags().BoolVar(&showTemplList, "list", false, "show all templ names")
}

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
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/algrvvv/ali/logger"
	"github.com/algrvvv/ali/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	errFailedToInitNewTempl = errors.New("failed to init new templ")
	errTemplAlreadyExists   = errors.New("template already exists\nuse --force for force rewrite")
)

var (
	createNewTempl bool
	forceAction    bool
	editTempl      bool
	showTemplList  bool
	templCmd       = &cobra.Command{
		Use:   "templ",
		Short: "use template configuration",
		Long:  `init new configuration by template`,
		Run: func(cmd *cobra.Command, args []string) {
			if showTemplList {
				if err := showListOfTemps(); err != nil {
					fmt.Println("failed to get list of templs")
					wrappedErr := fmt.Errorf("failed to show list of temps: %v", err)
					utils.CheckError(wrappedErr)
				}

				return
			}

			if len(args) != 1 {
				fmt.Println("failed to start `templ` command: expected templ name\nuse: ali templ templName --new")
				return
			}

			templName := args[0]
			logger.SaveDebugf("got templ name: %s", templName)

			if createNewTempl {
				if err := createNewTemplFunc(templName); err != nil {
					if errors.Is(err, errTemplAlreadyExists) {
						fmt.Println(err.Error())
						return
					}

					fmt.Println(errFailedToInitNewTempl.Error())
					// NOTE: оборачиываем ошибку для более удобного дальнейшего чтения
					// и либо просто говорим о том, что произошла ошибка, чекните логи,
					// либо выводим их
					wrappedErr := fmt.Errorf("failed to create new templ: %v", err)
					utils.CheckError(wrappedErr)
					return
				}
			}

			if editTempl {
				if err := editTemplByName(templName); err != nil {
					fmt.Println("failed to edit template")
					wrappedErr := fmt.Errorf("failed to edit templ: %v", err)
					utils.CheckError(wrappedErr)
				}

				return
			}

			if err := useTempl(templName); err != nil {
				fmt.Println("failed to use template")
			}
		},
	}
)

func useTempl(templName string) error {
	fmt.Println("use template: ", templName)

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %v", err)
	}

	templPath := templName + ".yml"
	path := filepath.Join(home, ".ali/", utils.TemplateDirName, templPath)
	logger.SaveDebugf("got path for template: %q", path)

	wordDir, err := os.Getwd()
	if err != nil {
		return err
	}

	localConfigFile := filepath.Join(wordDir, ".ali")
	_, err = os.Stat(localConfigFile)
	if err == nil {
		if !forceAction {
			fmt.Println("local config file already exists")
			fmt.Println("use --force for force rewrite")
			return nil
		}

		templViper := viper.New()
		templViper.SetConfigFile(path)
		templViper.SetConfigType("yaml")

		if err := templViper.ReadInConfig(); err != nil {
			return err
		}

		for _, key := range templViper.AllKeys() {
			value := templViper.Get(key)
			localViper.Set(key, value)
		}

		if err := localViper.WriteConfig(); err != nil {
			return err
		}

		fmt.Println("local config rewrited")
		return nil
	}

	fmt.Println("init new local config by template: ", templName)

	templConfig, err := os.Open(path)
	if err != nil {
		return err
	}
	defer templConfig.Close()

	localConfig, err := os.Create(localConfigFile)
	if err != nil {
		return err
	}
	defer localConfig.Close()

	_, err = io.Copy(localConfig, templConfig)
	if err != nil {
		return err
	}

	fmt.Println("local config from template created")
	return nil
}

func showListOfTemps() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %v", err)
	}

	path := filepath.Join(home, ".ali/", utils.TemplateDirName)
	logger.SaveDebugf("got path for templates: %q", path)

	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, entry := range dir {
		templName := utils.GetTemplNameByFile(entry.Name())
		fmt.Printf("%d. %s\n", i+1, templName)
	}

	return nil
}

func editTemplByName(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %v", err)
	}

	templPath := name + ".yml"
	path := filepath.Join(home, ".ali/", utils.TemplateDirName, templPath)
	logger.SaveDebugf("got path for template: %q", path)

	editor := viper.GetString("app.editor")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run template editor: %v", err)
	}

	return nil
}

func createNewTemplFunc(templName string) error {
	templPath := templName + ".yml"
	logger.SaveDebugf("init new template")
	fmt.Println("init new template by name: ", templName)

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %v", err)
	}

	path := filepath.Join(home, ".ali/", utils.TemplateDirName, templPath)
	logger.SaveDebugf("got path for template: %q", path)

	_, err = os.Stat(path)
	if err == nil {
		if forceAction {
			// if err := os.Remove(path); err != nil {
			// 	return fmt.Errorf("failed to remove old template: %v", err)
			// }

			// NOTE: удалять файл не обязательно, так как os.Create выглядит так:
			//  return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
			// как видно, функция использует O_TRUNC, а значит перезапишет файл
			fmt.Printf("old template by name %q truncated\n", templName)
		} else {
			return errTemplAlreadyExists
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create templ file: %v", err)
	}
	file.Close()

	fmt.Println("template file created")

	editor := viper.GetString("app.editor")
	if editor == "" {
		editor = "vi"
	}
	logger.SaveDebugf("got editor: %s", editor)

	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run template editor: %v", err)
	}

	fmt.Println()
	fmt.Println("template saved by name: ", templName)
	fmt.Printf("for run: ali templ %s\n", templName)
	fmt.Printf("for edit: ali templ %s --edit\n", templName)

	return nil
}

func init() {
	rootCmd.AddCommand(templCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	templCmd.Flags().BoolVar(&createNewTempl, "new", false, "create new template")
	templCmd.Flags().BoolVar(&forceAction, "force", false, "create new template force (rewrite existed templ)")
	templCmd.Flags().BoolVar(&editTempl, "edit", false, "edit template by name")
	templCmd.Flags().BoolVar(&showTemplList, "list", false, "show all templ names")
}

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
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/algrvvv/ali/logger"
	"github.com/algrvvv/ali/utils"
)

const (
	color      = "\033[1;34m"
	resetColor = "\033[0m"
)

// listCmd represents the list command
var (
	tablePrint     bool
	printVariables bool
	printEnvsFlag  bool

	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "Get list aliases",
		Example: "ali list - to see all aliases\nali list push - to search aliases that contains 'push' in name",
		Run: func(_ *cobra.Command, args []string) {
			var search string
			if len(args) == 1 {
				search = args[0]
			}

			if printVariables {
				printVars(search)
				return
			}

			if printEnvsFlag {
				printEnvs(search)
				return
			}

			printAliases(search)
		},
	}
)

func printEnvs(search string) {
	logger.SaveDebugf("print envs")
	fmt.Println("Available Envs:")

	envs := viper.GetStringMap("env")

	aliases := utils.LoadAliases(viper.GetViper())
	for alias, entry := range aliases {
		for name, value := range entry.Env {
			key := fmt.Sprintf("%s (%s)", name, alias)
			envs[key] = value
		}
	}

	if tablePrint {
		envsTablePrint(envs, search)
		return
	}

	envsFullPrint(envs, search)
}

func envsFullPrint(envs map[string]any, search string) {
	var count int
	for name, value := range envs {
		if !searchInEnvs(search, name) {
			continue
		}

		prefix := "  └── "
		if count != len(envs)-1 {
			prefix = "  ├── "
		}
		count++

		name = "$" + strings.ToUpper(name)
		name = utils.Colors["red"] + name + resetColor

		re := regexp.MustCompile(`\(([^)]+)\)`)
		name = re.ReplaceAllStringFunc(name, func(alias string) string {
			return color + strings.ToLower(alias) + resetColor
		})

		fmt.Printf("%s%s -> %v\n", prefix, name, value)
	}
}

func envsTablePrint(envs map[string]any, search string) {
	fmt.Printf("+%s+%s+\n", strings.Repeat("-", 30), strings.Repeat("-", 42))
	fmt.Printf("| Env%s| Command%s|\n",
		strings.Repeat(" ", 32-len(" alias")),
		strings.Repeat(" ", 42-len(" command")),
	)
	fmt.Printf("+%s+%s+\n", strings.Repeat("-", 30), strings.Repeat("-", 42))

	for name, value := range envs {
		if !searchInEnvs(search, name) {
			continue
		}

		name = "$" + strings.ToUpper(name)
		re := regexp.MustCompile(`\(([^)]+)\)`)
		name = re.ReplaceAllStringFunc(name, func(alias string) string {
			return strings.ToLower(alias)
		})

		fmt.Printf("| %s%-28s%s | %-40s |\n",
			utils.Colors["red"],
			utils.TruncateString(name, 28),
			resetColor,
			utils.TruncateString(fmt.Sprintf("%v", value), 40),
		)

		fmt.Printf("+%s+%s+\n", strings.Repeat("-", 30), strings.Repeat("-", 42))
	}
}

func searchInEnvs(search, envName string) bool {
	if search == "" {
		return true
	}

	if strings.Contains(envName, search) {
		return true
	}

	return false
}

func printVars(search string) {
	vars, err := utils.GetVars()
	if err != nil {
		fmt.Println("failed to get vars")
		logger.SaveDebugf("failed to get vars: %s", err)
		os.Exit(1)
	}

	logger.SaveDebugf("print variables")
	fmt.Println("Available Variables:")

	if tablePrint {
		varsTablePrint(vars, search)
		return
	}

	varsFullPrint(vars, search)
}

func varsFullPrint(vars map[string]string, search string) {
	var count int
	for name, value := range vars {
		if !searchInVars(search, name, value) {
			continue
		}

		prefix := "  └── "
		if count != len(vars)-1 {
			prefix = "  ├── "
		}
		count++

		fmt.Printf("%s%s%s%s -> %s\n", prefix, utils.Colors["lime"], name, resetColor, value)
	}
}

func varsTablePrint(vars map[string]string, search string) {
	fmt.Printf("+%s+%s+\n", strings.Repeat("-", 22), strings.Repeat("-", 42))
	fmt.Printf("| Var%s| Command%s|\n",
		strings.Repeat(" ", 24-len(" alias")),
		strings.Repeat(" ", 42-len(" command")),
	)
	fmt.Printf("+%s+%s+\n", strings.Repeat("-", 22), strings.Repeat("-", 42))

	for name, value := range vars {
		if !searchInVars(search, name, value) {
			continue
		}

		fmt.Printf("| %s%-20s%s | %-40s |\n",
			utils.Colors["lime"],
			utils.TruncateString(name, 20),
			resetColor,
			utils.TruncateString(value, 40),
		)

		fmt.Printf("+%s+%s+\n", strings.Repeat("-", 22), strings.Repeat("-", 42))
	}
}

func searchInVars(search, varName, varValue string) bool {
	if search == "" {
		return true
	}

	if strings.Contains(varName, search) {
		return true
	}

	if strings.Contains(varValue, search) {
		return true
	}

	return false
}

func printAliases(search string) {
	fmt.Println("Available Aliases:")
	aliases := utils.LoadAliases(viper.GetViper())

	if tablePrint {
		aliasTablePrint(aliases, search)
		return
	}
	aliasFullPrint(aliases, search)
}

func aliasFullPrint(aliases map[string]utils.AliasEntry, search string) {
	var count int
	for alias, entry := range aliases {
		if !searchInAlias(search, alias, entry) {
			continue
		}

		prefix := "  └── "
		if count != len(aliases)-1 {
			prefix = "  ├── "
		}
		count++

		if entry.Desc == "" {
			entry.Desc = "no desc"
		}
		entry.Desc = fmt.Sprintf("%s%s%s", utils.Colors["yellow"], entry.Desc, resetColor)

		clr := color
		if entry.Parallel {
			clr = utils.Colors["cyan"]
		}

		if len(entry.Aliases) > 0 {
			alias += fmt.Sprintf(" %s(%s)%s", utils.Colors["orange"], strings.Join(entry.Aliases, ", "), resetColor)
		}

		fmt.Printf("%s%s%s%s -> %s\n", prefix, clr, alias, resetColor, entry.Desc)

		for i, c := range entry.Cmds {
			prefix := "     └──"
			if i != len(entry.Cmds)-1 {
				prefix = "     ├──"
			}

			// добавляем выделение переменных (vars)
			re := regexp.MustCompile(`\{\{\w+\}\}`)
			c = re.ReplaceAllStringFunc(c, func(varStr string) string {
				return utils.Colors["lime"] + varStr + resetColor
			})

			// добавляем подсвечивание переменных окружения (env)
			re = regexp.MustCompile(`\$(\w+)`)
			c = re.ReplaceAllStringFunc(c, func(envStr string) string {
				return utils.Colors["red"] + envStr + resetColor
			})

			fmt.Printf("%s %s\n", prefix, c)
		}
	}
}

func aliasTablePrint(aliases map[string]utils.AliasEntry, search string) {
	fmt.Printf("+%s+%s+%s+\n", strings.Repeat("-", 22), strings.Repeat("-", 42), strings.Repeat("-", 30))
	fmt.Printf("| Alias%s| Command%s| Description%s|\n",
		strings.Repeat(" ", 22-len(" alias")),
		strings.Repeat(" ", 42-len(" command")),
		strings.Repeat(" ", 30-len(" description")),
	)
	fmt.Printf("+%s+%s+%s+\n", strings.Repeat("-", 22), strings.Repeat("-", 42), strings.Repeat("-", 30))

	for alias, entry := range aliases {
		if !searchInAlias(search, alias, entry) {
			continue
		}

		var cmd string
		if len(entry.Cmds) == 1 {
			cmd = entry.Cmds[0]
		} else {
			// join string using ';'
			cmd = strings.Join(entry.Cmds, "; ")
		}

		// показываем синонимы
		if len(entry.Aliases) > 0 {
			alias += " (" + strings.Join(entry.Aliases, ", ") + ")"
		}

		// if len(entry.Aliases) > 0 {
		// 	alias += fmt.Sprintf(" %s(%s)%s", utils.Colors["orange"], strings.Join(entry.Aliases, ", "), resetColor)
		// }

		// отдаем заглушку для пустого описания
		if entry.Desc == "" {
			entry.Desc = "no desc"
		}

		clr := color
		if entry.Parallel {
			clr = utils.Colors["cyan"]
		}

		fmt.Printf("| %s%-20s%s | %-40s | %-28s |\n",
			clr,
			utils.TruncateString(alias, 20),
			resetColor,
			utils.TruncateString(cmd, 40),
			utils.TruncateString(entry.Desc, 28),
		)

		fmt.Printf("+%s+%s+%s+\n", strings.Repeat("-", 22), strings.Repeat("-", 42), strings.Repeat("-", 30))
	}
}

func searchInAlias(search, alias string, entry utils.AliasEntry) bool {
	// пропуск поиска
	if search == "" {
		return true
	}

	// поиск по основнопу алиасу
	if strings.Contains(alias, search) {
		return true
	}

	// поиск по описанию
	if strings.Contains(entry.Desc, search) {
		return true
	}

	// поиск по синонимам
	for _, syn := range entry.Aliases {
		if strings.Contains(syn, search) {
			return true
		}
	}

	for _, cmd := range entry.Cmds {
		if strings.Contains(cmd, search) {
			return true
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().BoolVarP(&tablePrint, "table", "T", false, "use table print")
	listCmd.Flags().BoolVarP(&printVariables, "vars", "v", false, "print variables")
	listCmd.Flags().BoolVarP(&printEnvsFlag, "envs", "e", false, "print envs")
}

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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/algrvvv/ali/v2/aliases"
	"github.com/algrvvv/ali/v2/config"
	"github.com/algrvvv/ali/v2/envs"
	"github.com/algrvvv/ali/v2/exec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/algrvvv/ali/v2/logger"
	"github.com/algrvvv/ali/v2/utils"
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("failed to get home dir: ", err)
		os.Exit(1)
	}

	path := filepath.Join(home, ".ali")
	if err := os.Mkdir(path, 0777); err != nil && !errors.Is(err, os.ErrExist) {
		fmt.Println("failed to create global ali dir: ", err)
		os.Exit(1)
	}
}

var (
	localViper *viper.Viper

	// doParallel         bool
	debug              bool
	localEnv           bool
	withoutOutput      bool
	outputColor        string
	printResultCommand bool
	aliEnvFileFlag     string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:                "ali",
		Short:              "ali - cli app for your aliases",
		Args:               cobra.ArbitraryArgs,
		ValidArgsFunction:  aliases.GetValidAliases,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("use help for see usage")
				return
			}

			alias := args[0]
			//params := args[1:]
			unknownArgs, unknownFlags := utils.ParseUnknownArgs(os.Args[2:])
			params := strings.Split(unknownArgs, " ")

			logger.SaveDebugf("got alias: %s", alias)
			logger.SaveDebugf("got params(%d): %v", len(params), params)
			logger.SaveDebugf("got unknown args: %v", unknownArgs)
			logger.SaveDebugf("got unknown flags: %v", unknownFlags)

			als := aliases.Load(viper.GetViper())
			aliasEntry := aliases.SearchSynonyms(als, alias)
			logger.SaveDebugf("got alias entry: %v", aliasEntry)

			if aliasEntry == nil {
				fmt.Println("alias not found")
				return
			}

			env := envs.Load(aliasEntry, aliEnvFileFlag)

			if aliasEntry.Parallel {
				exec.ExecuteParallel(
					aliasEntry,
					params,
					unknownFlags,
					env,
					printResultCommand,
				)
			} else {
				for _, command := range aliasEntry.Cmds {
					err := exec.ExecuteLocal(
						command,
						aliasEntry.Dir,
						params,
						unknownFlags,
						env,
						printResultCommand,
					)
					if err != nil {
						fmt.Println("failed to get cmd: ", err)
						return
					}
				}
			}
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "print debug messages")
	rootCmd.PersistentFlags().BoolVarP(&localEnv, "local-env", "L", false, "use only local env")
	rootCmd.PersistentFlags().BoolVar(&withoutOutput, "without-output", false, "dont show parallel commands output")
	rootCmd.PersistentFlags().StringVar(&outputColor, "output-color", "", "color of the ouput of the parallel command")
	rootCmd.PersistentFlags().StringVar(&aliEnvFileFlag, "env-file", ".alienv", "file for load hidden envs")
	rootCmd.Flags().BoolVar(&printResultCommand, "print", false, "print result command before start exec")

	// NOTE: not used
	// rootCmd.PersistentFlags().BoolVarP(&doParallel, "parallel", "p", false, "do parallel command")

	// WARN: only for dev
	// rootCmd.PersistentFlags().StringVar(&localConfig, "local-config", ".ali", "local config path")
}

func initConfig() {
	if !localEnv {
		config.InitGlobal(rootCmd)
	}

	config.InitLocal(localViper, localEnv)
	config.InitIncluded()
}

func initLogger() {
	home, err := os.UserHomeDir()
	utils.CheckError(err)

	path := filepath.Join(home, ".ali/ali.log")
	err = logger.NewLogger(path, &logger.Options{Debug: true, MoreInfo: false, Stdout: debug})
	utils.CheckError(err)
}

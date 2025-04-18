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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/algrvvv/ali/logger"
	"github.com/algrvvv/ali/parallel"
	"github.com/algrvvv/ali/utils"
)

const localConfig = ".ali"

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
	debug              bool
	localEnv           bool
	doParallel         bool
	withoutOutput      bool
	outputColor        string
	printResultCommand bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:                "ali",
		Short:              "ali - cli app for your aliases",
		Args:               cobra.ArbitraryArgs,
		ValidArgsFunction:  getAliases,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("use help for see usage")
				return
			}
			alias := args[0]
			params := args[1:]
			unknownFlags := parseUnknownFlags(os.Args[1:])

			logger.SaveDebugf("got alias: %s", alias)
			logger.SaveDebugf("got params(%d): %v", len(params), params)
			logger.SaveDebugf("got unknown flags: %v", unknownFlags)

			if doParallel {
				parallel.DoParrallel(alias, outputColor, withoutOutput)
			} else {
				command := utils.GetAlias(alias)
				logger.SaveDebugf("got command: %s", command)

				utils.ExecuteAlias(command, params, unknownFlags, printResultCommand)
			}
		},
	}
)

func getAliases(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	aliases, ok := viper.Get("aliases").(map[string]any)
	if !ok {
		fmt.Println("failed to get all aliases")
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	var res []string
	for alias, command := range aliases {
		res = append(res, fmt.Sprintf("%s%s%s\t%s", color, alias, resetColor, command))
	}

	return res, cobra.ShellCompDirectiveNoFileComp
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger, initConfig, initLocalConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ali.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "print debug messages")
	rootCmd.PersistentFlags().BoolVarP(&localEnv, "local-env", "L", false, "use only local env")
	rootCmd.PersistentFlags().BoolVarP(&doParallel, "parallel", "p", false, "do parallel command")
	rootCmd.PersistentFlags().BoolVar(&withoutOutput, "without-output", false, "dont show parallel commands output")
	rootCmd.PersistentFlags().StringVar(&outputColor, "output-color", "", "color of the ouput of the parallel command")
	rootCmd.Flags().BoolVar(&printResultCommand, "print", false, "print result command before start exec")
}

func initConfig() {
	if !localEnv {
		initGlobalConfig()
	}

	initLocalConfig()
}

func initGlobalConfig() {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err != nil {
		fmt.Println("error occurred: ", err)
		os.Exit(1)
	}

	// пропускаем если setup скип
	if cmd.Name() == "setup" || cmd.Name() == "version" || cmd.Name() == "help" {
		return
	}

	home, err := os.UserHomeDir()
	utils.CheckError(err)

	path := filepath.Join(home, ".ali")
	viper.AddConfigPath(path)
	viper.SetConfigType("toml")
	viper.SetConfigName("config")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		utils.CheckError(err)
	}
	logger.SaveDebugf("using config: %s", viper.ConfigFileUsed())
}

func initLocalConfig() {
	dir, err := os.Getwd()
	utils.CheckError(err)
	path := filepath.Join(dir, localConfig)

	confType, err := utils.GetConfigurationType(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.SaveDebugf("failed to get configuration type: %v", err)
	}

	logger.SaveDebugf("read local config as %s", confType)

	viper.SetConfigName(".ali")
	viper.SetConfigType(confType)
	viper.AddConfigPath(".")

	if localEnv {
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			logger.SaveDebugf("load local config error: %v", err)
			// utils.CheckError(err)
		}
		logger.SaveDebugf("using config: %s", viper.ConfigFileUsed())
	} else {
		if err := viper.MergeInConfig(); err != nil {
			logger.SaveDebugf("local config not found")
		} else {
			logger.SaveDebugf("local config loaded")
		}
	}
}

func initLogger() {
	home, err := os.UserHomeDir()
	utils.CheckError(err)

	path := filepath.Join(home, ".ali/ali.log")
	err = logger.NewLogger(path, &logger.Options{Debug: true, MoreInfo: false, Stdout: debug})
	utils.CheckError(err)
}

func parseUnknownFlags(args []string) map[string]string {
	flags := make(map[string]string)
	for _, arg := range args {
		// NOTE: пропускаем зарезервированный ключ
		// TODO: придумать как избежать этого, если вдруг нужно будет использовать такой ключ
		if arg == "-D" || arg == "--debug" {
			continue
		}

		if strings.HasPrefix(arg, "-") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				flags[parts[0]] = parts[1]
			} else {
				flags[parts[0]] = ""
			}
		}
	}
	return flags
}

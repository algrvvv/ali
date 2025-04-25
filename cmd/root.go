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
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/algrvvv/ali/local"
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
	localViper *viper.Viper

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

			aliases := utils.LoadAliases(viper.GetViper())
			aliasEntry := utils.SearchSynonyms(aliases, alias)
			logger.SaveDebugf("got alias entry: %v", aliasEntry)

			if aliasEntry == nil {
				fmt.Println("alias not found")
				return
			}

			envs := utils.GetEnvs(aliasEntry)

			if aliasEntry.Parallel {
				parallel.ExecuteParallel(
					aliasEntry,
					params,
					unknownFlags,
					envs,
					printResultCommand,
				)
			} else {
				for _, command := range aliasEntry.Cmds {
					err := local.ExecuteLocal(
						command,
						aliasEntry.Dir,
						params,
						unknownFlags,
						envs,
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
	cobra.OnInitialize(initLogger, initConfig)

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

	// WARN: only for dev
	// rootCmd.PersistentFlags().StringVar(&localConfig, "local-config", ".ali", "local config path")
}

func initConfig() {
	if !localEnv {
		initGlobalConfig()
	}

	initLocalConfig()
	initInclideConfigs()
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
	viper.SetConfigType(utils.YamlConfigurationType)
	viper.SetConfigName("config")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		utils.CheckError(err)
	}
	logger.SaveDebugf("using config: %s", viper.ConfigFileUsed())
}

func initLocalConfig() {
	logger.SaveDebugf("local config: %s", localConfig)

	localViper = viper.New()
	localViper.SetConfigName(localConfig)
	localViper.SetConfigType(utils.YamlConfigurationType)
	localViper.AddConfigPath(".")

	viper.SetConfigName(localConfig)
	viper.SetConfigType(utils.YamlConfigurationType)
	viper.AddConfigPath(".")

	if err := localViper.ReadInConfig(); err != nil {
		logger.SaveDebugf("load local config error: %v", err)
		// utils.CheckError(err)
	} else {
		logger.SaveDebugf("local viper read config successfully")
	}

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

func initInclideConfigs() {
	include := viper.GetStringSlice("include")
	logger.SaveDebugf("includes: %v", include)

	for _, includePath := range include {
		if strings.HasPrefix(includePath, "~") {
			logger.SaveDebugf("user use ~ in dir param")
			home, err := os.UserHomeDir()
			if err != nil {
				logger.SaveDebugf("failed to get user home dir")
				return
			}

			logger.SaveDebugf("got home user dir: %q", home)
			includePath = strings.Replace(includePath, "~", home, 1)
			logger.SaveDebugf("command dir after change ~ to home dir: %q", includePath)
		}
		logger.SaveDebugf("include path: %s", includePath)

		info, err := os.Stat(includePath)
		if err != nil {
			fmt.Println("failed to get include file stat: ", err)
			logger.SaveDebugf("failed to get include file stat: %v", err)
			return
		}

		path := includePath
		if info.IsDir() {
			path = filepath.Join(path, ".ali")
		}

		logger.SaveDebugf("final unclude path: %s", path)
		viper.SetConfigFile(path)
		if err := viper.MergeInConfig(); err != nil {
			logger.SaveDebugf("uncluded config not found")
		} else {
			logger.SaveDebugf("uncluded config loaded")
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
	reservedFlags := []string{
		"-D", "-debug", "--debug",
		"-print", "-print",
		"-L", "--local-config",
		"-local-config",
	}

	flags := make(map[string]string)
	for _, arg := range args {
		// NOTE: пропускаем зарезервированный ключ
		// TODO: придумать как избежать этого, если вдруг нужно будет использовать такой ключ
		if strings.Contains(arg, "=") {
			prepArg := strings.SplitN(arg, "=", 2)[0]
			if slices.Contains(reservedFlags, prepArg) {
				logger.SaveDebugf("got reserved flag: %s; skip", arg)
				continue
			}
		}

		if slices.Contains(reservedFlags, arg) {
			logger.SaveDebugf("got reserved flag: %s; skip", arg)
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
